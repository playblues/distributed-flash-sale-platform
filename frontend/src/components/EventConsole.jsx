import React, { useRef, useEffect } from 'react';

export function EventConsole({ logs, onClearLogs }) {
  const containerRef = useRef(null);

  useEffect(() => {
    if (containerRef.current) {
      containerRef.current.scrollTop = containerRef.current.scrollHeight;
    }
  }, [logs]);

  const handleCopy = () => {
    const text = logs.map((l) => `[${l.time}] [${l.type.toUpperCase()}] ${l.message}`).join('\n');
    navigator.clipboard.writeText(text);
  };

  const getLogBadge = (type) => {
    if (type === 'success') {
      return <span className="text-emerald-500 font-bold">[OK]</span>;
    }
    if (type === 'error') {
      return <span className="text-rose-500 font-bold">[ERR]</span>;
    }
    if (type === 'warn') {
      return <span className="text-amber-500 font-bold">[WARN]</span>;
    }
    if (type === 'ws') {
      return <span className="text-cyan-400 font-bold">[KAFKA/WS]</span>;
    }
    return <span className="text-slate-600">[INFO]</span>;
  };

  const getLogTextColor = (type) => {
    if (type === 'success') return 'text-emerald-400 font-medium';
    if (type === 'error') return 'text-rose-400 font-medium';
    if (type === 'warn') return 'text-amber-400';
    if (type === 'ws') return 'text-cyan-300 font-medium';
    return 'text-slate-300';
  };

  return (
    <div className="bg-slate-900 border border-slate-800 rounded-2xl p-5 shadow-xl flex-1 flex flex-col min-h-[300px]">
      <div className="flex justify-between items-center pb-3 border-b border-slate-800/80 mb-3">
        <div className="flex items-center gap-2">
          <div className="flex gap-1.5">
            <span className="w-2.5 h-2.5 rounded-full bg-rose-500/80"></span>
            <span className="w-2.5 h-2.5 rounded-full bg-amber-500/80"></span>
            <span className="w-2.5 h-2.5 rounded-full bg-emerald-500/80"></span>
          </div>
          <span className="text-xs font-mono text-slate-400 ml-2">Event Stream Console</span>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={handleCopy}
            className="text-[11px] text-slate-500 hover:text-slate-300 font-mono transition px-2 py-1 rounded hover:bg-slate-800"
          >
            Копировать
          </button>
          <button
            onClick={onClearLogs}
            className="text-[11px] text-slate-500 hover:text-slate-300 font-mono transition px-2 py-1 rounded hover:bg-slate-800"
          >
            Очистить
          </button>
        </div>
      </div>

      <div
        ref={containerRef}
        className="flex-1 bg-slate-950 border border-slate-800/60 rounded-xl p-3.5 overflow-y-auto max-h-80 font-mono text-xs space-y-1.5"
      >
        {logs.length === 0 ? (
          <div className="text-slate-600">Логи пусты...</div>
        ) : (
          logs.map((item, idx) => (
            <div key={idx} className={`${getLogTextColor(item.type)} flex items-start gap-1.5 leading-relaxed`}>
              <span className="text-slate-600 text-[10px] select-none pt-0.5">{item.time}</span>
              {getLogBadge(item.type)}
              <span>{item.message}</span>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
