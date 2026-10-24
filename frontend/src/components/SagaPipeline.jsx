import React from 'react';

const STEPS = [
  { id: 'gateway', title: '1. Gateway', sub: 'RateLimit & JWT' },
  { id: 'inventory', title: '2. Inventory', sub: 'Redis Lua DECR' },
  { id: 'order', title: '3. Order Service', sub: 'Postgres & Outbox' },
  { id: 'kafka', title: '4. Kafka', sub: 'orders-topic' },
  { id: 'payment', title: '5. Payment', sub: 'Saga Decision' },
  { id: 'ws', title: '6. WebSocket', sub: 'Push в браузер' },
];

export function SagaPipeline({ pipelineState, pipelineText }) {
  const getStepClass = (stepId) => {
    const state = pipelineState[stepId];
    if (state === 'active') {
      return 'border-indigo-500 bg-indigo-950/40 text-indigo-300 font-bold shadow-lg shadow-indigo-500/20';
    }
    if (state === 'success') {
      return 'border-emerald-500 bg-emerald-950/40 text-emerald-300 font-bold';
    }
    if (state === 'failed') {
      return 'border-rose-500 bg-rose-950/40 text-rose-300 font-bold';
    }
    return 'border-slate-800 bg-slate-950/60 text-slate-400';
  };

  return (
    <section className="bg-slate-900/90 border border-slate-800/80 rounded-2xl p-5 shadow-2xl backdrop-blur">
      <div className="flex justify-between items-center mb-4">
        <span className="text-xs font-semibold uppercase tracking-wider text-slate-400 flex items-center gap-2">
          <svg className="w-4 h-4 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
          Saga Choreography Pipeline
        </span>
        <span className="text-xs font-mono text-slate-400">{pipelineText}</span>
      </div>

      <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-2 text-center text-xs font-mono">
        {STEPS.map((s) => (
          <div
            key={s.id}
            className={`p-3 rounded-xl border transition-all duration-300 ${getStepClass(s.id)}`}
          >
            <div className="font-semibold mb-0.5">{s.title}</div>
            <div className="text-[10px] text-slate-500">{s.sub}</div>
          </div>
        ))}
      </div>
    </section>
  );
}
