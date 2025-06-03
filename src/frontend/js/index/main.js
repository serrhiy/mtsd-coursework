'use strict';

import WebSocketTransport from "../WebSocketTransport.js";
import structure from "./structure.js";
import scaffold from "../scaffold.js";
import { getRandomUsername, debounce } from "./utils.js";

const saveChatButton = document.getElementById('createRoomButton');
const usernameInput = document.getElementById('usernameInput');
const roomNameInput = document.getElementById('roomNameInput');

const onUsernameChanged = async (api) => {
  const username = usernameInput.value.trim();
  if (username.length < 3 || username.length > 64) return;
  const token = localStorage.getItem('token');
  try {
    await api.users.update({ username, token });
    localStorage.setItem('username', username);
  } catch (error) {
    console.error(error);
  }
};

const onChatButton = async (api) => {
  const room = roomNameInput.value.trim();
  if (room.length < 3 || room.length > 64) return;
  const token = localStorage.getItem('token');
  try {
    await api.rooms.create({ room, token });
  } catch (error) {
    console.error(error);
  }
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
  const usernameChanged = onUsernameChanged.bind(null, api);
  usernameInput.addEventListener('input', debounce(3000, usernameChanged));
  saveChatButton.addEventListener('click', onChatButton.bind(null, api));
};

main();
