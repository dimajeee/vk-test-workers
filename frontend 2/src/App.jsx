// frontend/src/App.jsx
import React, { useState, useEffect } from "react";

function App() {
  const [workers, setWorkers] = useState(0);
  const [queueLength, setQueueLength] = useState(0);
  const [messagesProcessed, setMessagesProcessed] = useState(0);
  const [messagesTotal, setMessagesTotal] = useState(0);

  const [addCount, setAddCount] = useState(1);
  const [removeCount, setRemoveCount] = useState(1);
  const [sendCount, setSendCount] = useState(10);
  const [log, setLog] = useState([]);

  const apiBase = "http://127.0.0.1:8080";

  const addLog = (message) => {
    setLog((logs) => [message, ...logs].slice(0, 100));
  };

  const fetchStats = async () => {
    try {
      const res = await fetch(`${apiBase}/stats`);
      if (!res.ok) throw new Error("Failed to fetch stats");
      const data = await res.json();

      setWorkers(data.workers);
      setQueueLength(data.queue_length);
      setMessagesProcessed(data.messages_processed);
      setMessagesTotal(data.messages_total);
    } catch (e) {
      addLog("Failed to fetch stats");
    }
  };

  useEffect(() => {
    fetchStats();
    const interval = setInterval(fetchStats, 2000);
    return () => clearInterval(interval);
  }, []);

  const addWorkers = async () => {
    try {
      const res = await fetch(`${apiBase}/workers/add/${addCount}`, { method: "POST" });
      if (res.ok) {
        addLog(`Added ${addCount} workers`);
        fetchStats();
      } else {
        addLog("Failed to add workers");
      }
    } catch {
      addLog("Error adding workers");
    }
  };

  const removeWorkers = async () => {
    try {
      const res = await fetch(`${apiBase}/workers/remove/${removeCount}`, { method: "POST" });
      if (res.ok) {
        addLog(`Requested removal of ${removeCount} workers`);
        fetchStats();
      } else {
        addLog("Failed to remove workers");
      }
    } catch {
      addLog("Error removing workers");
    }
  };

  const sendMessages = async () => {
    try {
      const messages = Array.from({ length: sendCount }, (_, i) => `msg-${i + 1}`);
      const res = await fetch(`${apiBase}/send`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ messages }),
      });
      if (res.ok) {
        addLog(`Sent ${sendCount} messages`);
        fetchStats();
      } else {
        addLog("Failed to send messages");
      }
    } catch {
      addLog("Error sending messages");
    }
  };

  const stopAll = async () => {
    try {
      const res = await fetch(`${apiBase}/stop`, { method: "POST" });
      if (res.ok) {
        addLog("Stopped all workers");
        fetchStats();
      } else {
        addLog("Failed to stop all workers");
      }
    } catch {
      addLog("Error stopping all workers");
    }
  };

  // Рассчитываем прогресс обработки сообщений
  const progressPercent =
    messagesTotal > 0 ? Math.round((messagesProcessed / messagesTotal) * 100) : 0;

  return (
    <div style={{ padding: 20, fontFamily: "Arial" }}>
      <h1>Worker Pool Control Panel</h1>

      <div>
        <strong>Current Workers:</strong> {workers} <br />
        <strong>Queue Length:</strong> {queueLength} <br />
        <strong>Messages Processed:</strong> {messagesProcessed} / {messagesTotal}{" "}
        ({progressPercent}%)
      </div>

      <div
        style={{
          marginTop: 10,
          width: "100%",
          height: 20,
          backgroundColor: "#ddd",
          borderRadius: 10,
          overflow: "hidden",
        }}
      >
        <div
          style={{
            width: `${progressPercent}%`,
            height: "100%",
            backgroundColor: "#4caf50",
            transition: "width 0.5s ease-in-out",
          }}
        />
      </div>

      <div style={{ marginTop: 20 }}>
        <label>
          Add Workers:{" "}
          <input
            type="number"
            value={addCount}
            min={1}
            onChange={(e) => setAddCount(Number(e.target.value))}
            style={{ width: 60 }}
          />
        </label>
        <button onClick={addWorkers} style={{ marginLeft: 10 }}>
          Add
        </button>
      </div>

      <div style={{ marginTop: 10 }}>
        <label>
          Remove Workers:{" "}
          <input
            type="number"
            value={removeCount}
            min={1}
            onChange={(e) => setRemoveCount(Number(e.target.value))}
            style={{ width: 60 }}
          />
        </label>
        <button onClick={removeWorkers} style={{ marginLeft: 10 }}>
          Remove
        </button>
      </div>

      <div style={{ marginTop: 10 }}>
        <label>
          Send Messages:{" "}
          <input
            type="number"
            value={sendCount}
            min={1}
            onChange={(e) => setSendCount(Number(e.target.value))}
            style={{ width: 60 }}
          />
        </label>
        <button onClick={sendMessages} style={{ marginLeft: 10 }}>
          Send
        </button>
      </div>

      <div style={{ marginTop: 10 }}>
        <button
          onClick={stopAll}
          style={{ backgroundColor: "#d9534f", color: "white", padding: "6px 12px" }}
        >
          Stop All Workers
        </button>
      </div>

      <h3 style={{ marginTop: 30 }}>Logs</h3>
      <div
        style={{
          height: 200,
          overflowY: "scroll",
          backgroundColor: "#f0f0f0",
          padding: 10,
          borderRadius: 5,
          fontFamily: "monospace",
          fontSize: 12,
        }}
      >
        {log.length === 0 && <div>No logs yet</div>}
        {log.map((entry, i) => (
          <div key={i}>{entry}</div>
        ))}
      </div>
    </div>
  );
}

export default App;
