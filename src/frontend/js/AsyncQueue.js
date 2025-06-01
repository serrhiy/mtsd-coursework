'use strict';

export default class AsyncQueue {
  #futures = [];
  #values = [];

  put(value) {
    if (this.#futures.length > 0) {
      const future = this.#futures.shift();
      return void future(value);
    }
    this.#values.push(value);
  }

  get() {
    if (this.#values.length > 0) {
      return Promise.resolve(this.#values.shift());
    }
    return new Promise((resolve) => {
      this.#futures.push(resolve);
    });
  }

  [Symbol.asyncIterator]() {
    const next = async () => ({ value: await this.get(), done: false });
    return { next };
  }
}
