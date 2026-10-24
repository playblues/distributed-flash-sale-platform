import React from 'react';

export function StockCard({
  stock,
  initialStockTotal,
  userId,
  setUserId,
  price,
  setPrice,
  onRefreshStock,
  onInitStock,
  onBuy,
  onSpurt,
  isLoading,
  isSpurtLoading,
}) {
  const stockPercentage =
    stock !== null && initialStockTotal > 0
      ? Math.max(0, Math.min(100, Math.round((stock / initialStockTotal) * 100)))
      : 0;

  const getBarColor = () => {
    if (stock === 0) return 'bg-rose-500';
    if (stock !== null && stock <= 3) return 'bg-amber-500';
    return 'bg-gradient-to-r from-indigo-500 to-emerald-400';
  };

  const handleRandomUser = () => {
    const rnd = Math.floor(Math.random() * 900) + 100;
    setUserId(rnd);
  };

  return (
    <div className="bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-xl relative overflow-hidden">
      <div className="absolute top-0 right-0 w-32 h-32 bg-indigo-500/10 rounded-full blur-2xl -mr-10 -mt-10 pointer-events-none"></div>

      <div className="flex justify-between items-start">
        <div>
          <span className="text-[10px] uppercase font-bold tracking-wider text-indigo-400 bg-indigo-950/60 px-2.5 py-1 rounded-md border border-indigo-800/40">
            Flash Sale Event
          </span>
          <h2 className="text-xl font-bold text-white mt-2">Мировое турне: Рок-фестиваль</h2>
        </div>
        <span className="text-xs px-2.5 py-1 bg-emerald-950/80 border border-emerald-800/60 text-emerald-300 rounded-md font-mono">
          Event #1
        </span>
      </div>

      <div className="mt-5 p-3.5 bg-slate-950/80 border border-slate-800 rounded-xl space-y-2">
        <div className="flex justify-between text-xs">
          <span className="text-slate-400">Остаток билетов:</span>
          <span className="font-mono font-bold text-slate-200">
            {stock !== null ? `${stock} в наличии` : 'Загрузка...'}
          </span>
        </div>
        <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
          <div
            className={`h-full transition-all duration-500 ${getBarColor()}`}
            style={{ width: `${stockPercentage}%` }}
          ></div>
        </div>
        <div className="flex justify-between items-center pt-1">
          <button
            onClick={onRefreshStock}
            className="text-[11px] text-indigo-400 hover:text-indigo-300 flex items-center gap-1 transition"
          >
            <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Обновить остаток
          </button>
          <div className="flex items-center gap-1 text-[11px] text-slate-400">
            <span>Завоз:</span>
            {[5, 10, 50].map((amt) => (
              <button
                key={amt}
                onClick={() => onInitStock(amt)}
                className="px-2 py-0.5 bg-slate-800 hover:bg-slate-700 rounded text-slate-200 font-mono transition"
              >
                +{amt}
              </button>
            ))}
          </div>
        </div>
      </div>

      <div className="mt-5 space-y-4">
        <div>
          <div className="flex justify-between items-center mb-1">
            <label className="text-xs font-medium text-slate-300">ID пользователя (User ID):</label>
            <div className="flex gap-1">
              {[42, 101, 102].map((id) => (
                <button
                  key={id}
                  onClick={() => setUserId(id)}
                  className={`text-[10px] px-2 py-0.5 border rounded font-mono transition ${
                    userId === id
                      ? 'bg-indigo-600/40 border-indigo-500 text-indigo-200'
                      : 'bg-slate-800 border-slate-700 text-slate-300 hover:bg-slate-700'
                  }`}
                >
                  #{id}
                </button>
              ))}
              <button
                onClick={handleRandomUser}
                className="text-[10px] px-2 py-0.5 bg-slate-800 hover:bg-slate-700 border border-slate-700 rounded text-slate-300 transition"
                title="Случайный ID"
              >
                🎲
              </button>
            </div>
          </div>
          <input
            type="number"
            value={userId}
            onChange={(e) => setUserId(parseInt(e.target.value, 10) || 1)}
            className="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2.5 text-sm font-mono focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition"
          />
        </div>

        <div>
          <label className="text-xs font-medium text-slate-300 block mb-1">Сумма заказа (₽):</label>
          <input
            type="number"
            value={price}
            onChange={(e) => setPrice(parseFloat(e.target.value) || 0)}
            className="w-full bg-slate-950 border border-slate-700 rounded-xl px-3.5 py-2.5 text-sm font-mono focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 transition"
          />
        </div>
      </div>

      <div className="mt-5 space-y-2">
        <span className="text-[11px] font-semibold uppercase tracking-wider text-slate-400 block">
          Быстрые сценарии:
        </span>
        <div className="grid grid-cols-2 gap-2 text-xs">
          <button
            onClick={() => {
              setPrice(5000);
              onBuy(userId, 5000);
            }}
            disabled={isLoading}
            className="p-2.5 bg-slate-950 hover:bg-slate-800 border border-slate-800 hover:border-emerald-500/50 rounded-xl text-left transition flex flex-col disabled:opacity-50"
          >
            <span className="font-medium text-emerald-400">🟢 Успех (5 000 ₽)</span>
            <span className="text-[10px] text-slate-500">PENDING ➔ PAID</span>
          </button>

          <button
            onClick={() => {
              setPrice(150000);
              onBuy(userId, 150000);
            }}
            disabled={isLoading}
            className="p-2.5 bg-slate-950 hover:bg-slate-800 border border-slate-800 hover:border-amber-500/50 rounded-xl text-left transition flex flex-col disabled:opacity-50"
          >
            <span className="font-medium text-amber-400">🟡 Saga Rollback</span>
            <span className="text-[10px] text-slate-500">&gt; 100k ₽ ➔ CANCELLED</span>
          </button>

          <button
            onClick={() => {
              setPrice(666);
              onBuy(userId, 666);
            }}
            disabled={isLoading}
            className="p-2.5 bg-slate-950 hover:bg-slate-800 border border-slate-800 hover:border-rose-500/50 rounded-xl text-left transition flex flex-col disabled:opacity-50"
          >
            <span className="font-medium text-rose-400">🟣 Poison Pill (666 ₽)</span>
            <span className="text-[10px] text-slate-500">3 ретрая ➔ Kafka DLT</span>
          </button>

          <button
            onClick={onSpurt}
            disabled={isSpurtLoading || isLoading}
            className="p-2.5 bg-slate-950 hover:bg-slate-800 border border-slate-800 hover:border-indigo-500/50 rounded-xl text-left transition flex flex-col disabled:opacity-50"
          >
            <span className="font-medium text-indigo-400">⚡ Залп (5 покупок)</span>
            <span className="text-[10px] text-slate-500">
              {isSpurtLoading ? 'Отправка...' : 'Параллельный наплыв'}
            </span>
          </button>
        </div>
      </div>

      <div className="mt-6">
        <button
          onClick={() => onBuy(userId, price)}
          disabled={isLoading}
          className="w-full bg-indigo-600 hover:bg-indigo-500 text-white font-semibold py-3.5 px-4 rounded-xl transition duration-200 shadow-xl shadow-indigo-600/25 flex items-center justify-center gap-2 active:scale-[0.99] disabled:opacity-50"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
          </svg>
          <span>{isLoading ? 'Обработка запроса...' : 'Забронировать и купить'}</span>
        </button>
      </div>
    </div>
  );
}
