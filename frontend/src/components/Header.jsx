import React from 'react';

export function Header({ isConnected, userId, stock }) {
  return (
    <header className="w-full max-w-6xl mx-auto flex flex-col sm:flex-row justify-between items-start sm:items-center pb-6 border-b border-slate-800/80 mb-6 gap-4">
      <div className="flex items-center space-x-3">
        <div className="relative flex items-center justify-center">
          <span className="w-3.5 h-3.5 rounded-full bg-emerald-500"></span>
          <span className="absolute w-5 h-5 rounded-full bg-emerald-500/40 animate-ping"></span>
        </div>
        <div>
          <h1 className="text-xl font-bold tracking-tight text-white flex items-center gap-2">
            Flash-Sale Ticket Engine
            <span className="text-[10px] font-mono uppercase bg-indigo-950/80 text-indigo-400 border border-indigo-800/60 px-2 py-0.5 rounded">
              React + Polyglot
            </span>
          </h1>
          <p className="text-xs text-slate-400">
            Go Fiber + Spring Boot 3 + Redis Lua + Kafka KRaft + WebSockets
          </p>
        </div>
      </div>

      <div className="flex items-center flex-wrap gap-2 text-xs">
        <div
          className={`flex items-center gap-1.5 px-3 py-1.5 border rounded-lg font-mono transition-all ${
            isConnected
              ? 'bg-emerald-950/60 border-emerald-800 text-emerald-300'
              : 'bg-rose-950/50 border-rose-800/80 text-rose-300'
          }`}
        >
          <span className={`w-2 h-2 rounded-full ${isConnected ? 'bg-emerald-500' : 'bg-rose-500'}`}></span>
          <span>{isConnected ? `WebSocket: Online (User #${userId})` : 'WebSocket: Offline'}</span>
        </div>

        <div className="flex items-center gap-1.5 px-3 py-1.5 bg-slate-900 border border-slate-800 text-slate-300 rounded-lg font-mono">
          <span>📦 В наличии:</span>
          <span className="font-bold text-amber-400">
            {stock !== null ? `${stock} шт.` : '...'}
          </span>
        </div>
      </div>
    </header>
  );
}
