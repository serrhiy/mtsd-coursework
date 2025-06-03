'use strict';

import WebSocketTransport from "../WebSocketTransport.js";
import structure from "./structure.js";
import scaffold from "../scaffold.js";
import { getRandomUsername } from "./utils.js";

const saveUsernameButton = document.getElementById('saveUsernameButton');
const saveChatButton = document.getElementById('createRoomButton');
const usernameInput = document.getElementById('usernameInput');
const roomNameInput = document.getElementById('roomNameInput');

const onUsernameButton = (api) => {
  const source = usernameInput.value.trim();
  if (source.length === 0) return;
};

const onChatButton = (api) => {
  const source = roomNameInput.value.trim();
  if (source.length === 0) return;
};

const setupUsername = async () => {
  if (localStorage.getItem('username') === null) {
    const username = await getRandomUsername();
    localStorage.setItem('username', username);
  }
  const username = localStorage.getItem('username');
  usernameInput.value = username;
  return username;
}

const setupToken = async (api) => {
  if (localStorage.getItem('token') === null) {
    if (localStorage.getItem('username') === null) setupUsername();
    const username = localStorage.getItem('username');
    const token = await api.users.create(username);
    localStorage.setItem('token', token);
    return token;
  }
  const token = localStorage.getItem('token');
  const exists = await api.users.exists(token);
  if (exists) return token;
  localStorage.removeItem('token');
  return setupToken();
};

const main = async () => {
  const websocket = await new WebSocketTransport('ws://127.0.0.1:8080');
  const api = scaffold(structure, websocket);
  await setupUsername();
  await setupToken(api);
  saveUsernameButton.addEventListener('click', onUsernameButton.bind(null, api));
  saveChatButton.addEventListener('click', onChatButton.bind(null, api));
};

main();