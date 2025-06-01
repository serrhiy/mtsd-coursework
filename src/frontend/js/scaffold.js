'use strict';

export default (structure, transport) => {
  const api = Object.create(null);
  for (const [service, methods] of Object.entries(structure)) {
    const functions = {};
    for (const method of methods) {
      functions[method] = (data) => {
        return transport.send(service, method, data)
      };
    }
    api[service] = functions;
  }
  return api;
};
