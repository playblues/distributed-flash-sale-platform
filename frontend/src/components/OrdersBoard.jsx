import React from 'react';

export function OrdersBoard({ orders }) {
  const getBadgeClass = (status) => {
    if (status === 'PAID') {
      return 'bg-emerald-950/80 border-emerald-800 text-emerald-300';
    }
    if (status === 'CANCELLED') {
      return 'bg-rose-950/80 border-rose-800 text-rose-300';
    }
    return 'bg-amber-950/60 border-amber-800/80 text-amber-300 animate-pulse';
  };

  return (
    <div className="bg-slate-900 border border-slate-800 rounded-2xl p-5 shadow-xl">
      <div className="flex justify-between items-center mb-3">
        <h3 className="text-xs font-semibold uppercase tracking-wider text-slate-400 flex items-center gap-2">
          <svg className="w-4 h-4 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
          </svg>
          Созданные заказы (Живой трекер статусов)
        </h3>
        <span className="text-[10px] font-mono px-2 py-0.5 bg-slate-800 text-slate-400 rounded">
          {orders.length} заказов
        </span>
      </div>

      <div className="space-y-2 max-h-48 overflow-y-auto pr-1">
        {orders.length === 0 ? (
          <div className="text-center py-6 text-slate-600 text-xs">
            Заказов пока нет. Нажмите «Забронировать и купить» для запуска.
          </div>
        ) : (
          orders.map((order) => (
            <div
              key={order.id}
              className="p-3 bg-slate-950/80 border border-slate-800 rounded-xl flex items-center justify-between text-xs transition-all"
            >
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-lg bg-indigo-950/80 border border-indigo-800/60 flex items-center justify-center font-mono font-bold text-indigo-300">
                  #{order.id}
                </div>
                <div>
                  <div className="font-medium text-white flex items-center gap-2">
                    Заказ №{order.id}
                    <span className="text-[10px] text-slate-400 font-mono">User #{order.userId}</span>
                  </div>
                  <div className="text-[11px] font-mono text-slate-400 mt-0.5">
                    Сумма: <span className="text-slate-200 font-bold">{order.price.toLocaleString('ru-RU')} ₽</span>
                  </div>
                </div>
              </div>
              <div className="flex items-center gap-2">
                <span className={`px-2.5 py-1 rounded-md text-[10px] font-mono font-bold uppercase tracking-wider border ${getBadgeClass(order.status)}`}>
                  {order.status === 'PAID' ? 'PAID (Оплачен)' : order.status === 'CANCELLED' ? 'CANCELLED (Откат)' : order.status}
                </span>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
