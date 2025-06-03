'use strict';

const API = 'https://randomuser.me/api?results=1&nat=ua&inc=name';

export const getRandomUsername = async () => {
  const signal = AbortSignal.timeout(500);
  const response = await fetch(API, { signal });
  if (!response.ok) throw new Error('Invalid request');
  const json = await response.json();
  if (json.error) throw new Error(json.error);
  const { first, last } = json.results[0].name;
  return first + '-' + last;
};

export const debounce = (timeout, onTimeout) => {
  let timer = null;
  return () => {
    clearTimeout(timer);
    timer = setTimeout(onTimeout, timeout);
  }
};
