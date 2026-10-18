const WS_BASE_URL = 'ws://localhost:3002';

export class NotificationSocket {
  constructor(userId, onMessage, onStatusChange) {
    this.userId = userId;
    this.onMessage = onMessage;
    this.onStatusChange = onStatusChange;
    this.socket = null;
  }

  connect() {
    if (this.socket) {
      this.socket.close();
    }

    this.socket = new WebSocket(`${WS_BASE_URL}/ws/${this.userId}`);

    this.socket.onopen = () => {
      this.onStatusChange(true, this.userId);
    };

    this.socket.onmessage = (event) => {
      this.onMessage(event.data);
    };

    this.socket.onclose = () => {
      this.onStatusChange(false, this.userId);
    };

    this.socket.onerror = (err) => {
      console.error('WebSocket Error:', err);
    };
  }

  disconnect() {
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
  }
}