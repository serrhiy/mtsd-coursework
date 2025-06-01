'use strict';

import AsyncQueue from './AsyncQueue.js';

export default class WebSocketTransport extends WebSocket {
  #id = 0;
  #timeout = 0;
  #futures = new Map();
  #values = new AsyncQueue();

  constructor(url, timeout = 5000) {
    super(url);
    this.#timeout = timeout;
    this.addEventListener('message', (event) => {
      const payload = JSON.parse(event.data);
      if (payload.type === 'response' && this.#futures.has(payload.id)) {
        const { resolve, timer } = this.#futures.get(payload.id);
        this.#futures.delete(payload.id);
        resolve(payload.data);
        return void clearTimeout(timer);
      }
      this.#values.put(payload.data);
    });
    return new Promise((resolve, reject) => {
      const onOpen = () => {
        this.removeEventListener('error', reject);
        resolve(this);
      };
      this.addEventListener('open', onOpen, { once: true });
      this.addEventListener('error', reject, { once: true });
    });
  }

  send(service, method, data) {
    const payload = { id: this.#id, service, method, data };
    super.send(JSON.stringify(payload));
    return new Promise((resolve, reject) => {
      const onTimeout = () => {
        reject(new Error('Response waiting time is up'));
        this.#futures.delete(this.#id);
      };
      const timer = setTimeout(onTimeout, this.#timeout);
      this.#futures.set(this.#id, { resolve, timer });
      this.#id++;
    });
  }

  [Symbol.asyncIterator]() {
    return this.#values[Symbol.asyncIterator]();
  }
}
