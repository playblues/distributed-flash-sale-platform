import React, { useState, useEffect, useCallback } from 'react';
import { Header } from './components/Header';
import { SagaPipeline } from './components/SagaPipeline';
import { StockCard } from './components/StockCard';
import { OrdersBoard } from './components/OrdersBoard';
import { EventConsole } from './components/EventConsole';
import { login, buyTicket, initStock, getStock } from './api/client';
import { useWebSocket } from './hooks/useWebSocket';

const INITIAL_PIPELINE = {
  gateway: 'idle',
  inventory: 'idle',
  order: 'idle',
  kafka: 'idle',
  payment: 'idle',
  ws: 'idle',
};

export default function App() {
  const [authToken, setAuthToken] = useState('');
  const [userId, setUserId] = useState(42);
  const [price, setPrice] = useState(5000);
  const [stock, setStock] = useState(null);
  const [initialStockTotal, setInitialStockTotal] = useState(10);
  const [orders, setOrders] = useState([]);
  const [logs, setLogs] = useState([]);
  const [pipelineState, setPipelineState] = useState(INITIAL_PIPELINE);
  const [pipelineText, setPipelineText] = useState('Ожидание запуска запроса');
  const [isLoading, setIsLoading] = useState(false);
  const [isSpurtLoading, setIsSpurtLoading] = useState(false);

  const addLog = useCallback((message, type = 'info') => {
    const time = new Date().toLocaleTimeString();
    setLogs((prev) => [...prev, { time, type, message }]);
  }, []);

  const refreshStock = useCallback(async () => {
    const data = await getStock(1);
    if (data && data.stock !== undefined) {
      setStock(data.stock);
      setInitialStockTotal((prev) => Math.max(prev, data.stock));
    }
  }, []);

  const handleWebSocketMessage = useCallback((message) => {
    addLog(message, 'ws');
    const match = message.match(/№(\d+).*?:\s*([A-Z_]+)/);
    if (!match) return;

    const orderId = match[1];
    const paymentStatus = match[2];

    setPipelineState((prev) => ({
      ...prev,
      kafka: 'success',
      ws: 'success',
      payment: paymentStatus === 'SUCCESS' ? 'success' : 'failed',
    }));

    if (paymentStatus === 'SUCCESS') {
      setPipelineText('Оплата подтверждена (PAID)');
    } else {
      setPipelineText('Откат брони (Saga Compensated)');
      setTimeout(refreshStock, 600);
    }

    setOrders((prev) =>
      prev.map((o) =>
        String(o.id) === String(orderId)
          ? { ...o, status: paymentStatus === 'SUCCESS' ? 'PAID' : 'CANCELLED' }
          : o
      )
    );
  }, [addLog, refreshStock]);

  const { isConnected } = useWebSocket(userId, handleWebSocketMessage);

  useEffect(() => {
    addLog('Инициализация React-дашборда...', 'info');
    login()
      .then((token) => {
        setAuthToken(token);
        addLog('JWT-токен успешно получен от API Gateway', 'success');
      })
      .catch((e) => {
        addLog('Ошибка авторизации через Gateway: ' + e.message, 'error');
      });

    refreshStock();
  }, [addLog, refreshStock]);

  const handleInitStock = async (amount) => {
    try {
      await initStock(amount, 1);
      setInitialStockTotal((prev) => Math.max(prev, amount));
      await refreshStock();
      addLog(`На склад завезено +${amount} билетов для Event #1`, 'success');
    } catch (e) {
      addLog(`Ошибка завоза билетов: ${e.message}`, 'error');
    }
  };

  const handleBuy = async (targetUserId, targetPrice) => {
    setIsLoading(true);
    setPipelineState({ ...INITIAL_PIPELINE, gateway: 'active' });
    setPipelineText('1. Проверка лимитов и идемпотентности на Gateway...');

    try {
      setPipelineState((prev) => ({ ...prev, inventory: 'active' }));
      setPipelineText('2. Атомарное бронирование в Redis (Lua-скрипт)...');

      const { status, data } = await buyTicket({
        userId: targetUserId,
        eventId: 1,
        price: targetPrice,
        token: authToken,
      });

      if (status === 201) {
        setPipelineState({
          gateway: 'success',
          inventory: 'success',
          order: 'success',
          kafka: 'active',
          payment: 'idle',
          ws: 'idle',
        });
        setPipelineText(`3. Заказ №${data.id} создан (PENDING). Отправлен в Kafka...`);
        addLog(`Заказ №${data.id} создан (PENDING). Сумма: ${targetPrice} ₽. Ожидание Kafka...`, 'info');

        setOrders((prev) => [
          { id: data.id, userId: targetUserId, price: targetPrice, status: 'PENDING' },
          ...prev,
        ]);
        await refreshStock();

        if (targetPrice === 666) {
          setPipelineState((prev) => ({ ...prev, payment: 'failed' }));
          setPipelineText('Poison Pill (666 ₽) ➔ 3 ретрая ➔ Kafka DLT');
          addLog(`Заказ №${data.id} вызвал Poison Pill (уходит в orders-topic.DLT)`, 'warn');
        }
      } else if (status === 409) {
        setPipelineState((prev) => ({ ...prev, inventory: 'failed' }));
        const errText =
          data.error ||
          (data.status === 'already_reserved'
            ? 'У вас уже есть бронь!'
            : 'Билеты закончились');
        setPipelineText(`Отказ: ${errText}`);
        addLog(`Отказ Inventory (409 Conflict): ${errText}`, 'warn');
        await refreshStock();
      } else if (status === 429) {
        setPipelineState((prev) => ({ ...prev, gateway: 'failed' }));
        setPipelineText('Rate Limiter: 429 Too Many Requests');
        addLog('Сработал Rate Limiter шлюза: Слишком много запросов', 'warn');
      } else {
        setPipelineState((prev) => ({ ...prev, gateway: 'failed' }));
        setPipelineText(`Ошибка сервера (${status})`);
        addLog(`Ошибка сервера (${status}): ${JSON.stringify(data)}`, 'error');
        await refreshStock();
      }
    } catch (e) {
      setPipelineState((prev) => ({ ...prev, gateway: 'failed' }));
      setPipelineText('Сетевой сбой при обращении к Gateway');
      addLog(`Ошибка сети: ${e.message}`, 'error');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSpurt = async () => {
    setIsSpurtLoading(true);
    addLog('⚡ Запуск залпа из 5 параллельных запросов от разных пользователей...', 'info');

    const baseUid = Math.floor(Math.random() * 800) + 100;
    const promises = [1, 2, 3, 4, 5].map((i) => {
      const uid = baseUid + i;
      return buyTicket({
        userId: uid,
        eventId: 1,
        price: 5000,
        token: authToken,
      }).then((r) => {
        addLog(
          `Пользователь #${uid}: статус ${r.status} (${JSON.stringify(r.data)})`,
          r.status === 201 ? 'success' : 'warn'
        );
        if (r.status === 201) {
          setOrders((prev) => [
            { id: r.data.id, userId: uid, price: 5000, status: 'PENDING' },
            ...prev,
          ]);
        }
        return r;
      });
    });

    await Promise.allSettled(promises);
    await refreshStock();
    setIsSpurtLoading(false);
    addLog('⚡ Залп завершен. Проверьте остаток билетов!', 'success');
  };

  return (
    <div className="min-h-screen p-4 md:p-8 flex flex-col font-sans selection:bg-indigo-500 selection:text-white">
      <Header isConnected={isConnected} userId={userId} stock={stock} />

      <main className="w-full max-w-6xl mx-auto space-y-6 flex-1">
        <SagaPipeline pipelineState={pipelineState} pipelineText={pipelineText} />

        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
          <div className="lg:col-span-5 space-y-6">
            <StockCard
              stock={stock}
              initialStockTotal={initialStockTotal}
              userId={userId}
              setUserId={setUserId}
              price={price}
              setPrice={setPrice}
              onRefreshStock={refreshStock}
              onInitStock={handleInitStock}
              onBuy={handleBuy}
              onSpurt={handleSpurt}
              isLoading={isLoading}
              isSpurtLoading={isSpurtLoading}
            />
          </div>

          <div className="lg:col-span-7 space-y-6 flex flex-col">
            <OrdersBoard orders={orders} />
            <EventConsole logs={logs} onClearLogs={() => setLogs([])} />
          </div>
        </div>
      </main>

      <footer className="w-full max-w-6xl mx-auto pt-6 mt-8 border-t border-slate-800/60 text-center text-xs text-slate-500 flex flex-col sm:flex-row justify-between items-center gap-2">
        <div>Distributed Flash-Sale Engine &bull; React UI</div>
        <div className="flex gap-4 font-mono text-[11px]">
          <a href="http://localhost:8080" target="_blank" rel="noreferrer" className="hover:text-indigo-400 transition">
            Kafka AKHQ ↗
          </a>
          <a href="http://localhost:16686" target="_blank" rel="noreferrer" className="hover:text-indigo-400 transition">
            Jaeger Traces ↗
          </a>
          <a href="http://localhost:3000" target="_blank" rel="noreferrer" className="hover:text-indigo-400 transition">
            API Gateway ↗
          </a>
        </div>
      </footer>
    </div>
  );
}
