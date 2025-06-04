'use strict';

import WebSocketTransport from "../WebSocketTransport.js";
import structure from "./structure.js";
import scaffold from "../scaffold.js";
import { getRandomUsername, debounce } from "./utils.js";

const DOMAIN = 'localhost';

const saveChatButton = document.getElementById('createRoomButton');
const usernameInput = document.getElementById('usernameInput');
const roomNameInput = document.getElementById('roomNameInput');
const roomList = document.getElementById('roomList');

const buildRoomLink = (title, token, username) => {
  const li = document.createElement('li');
  const a = document.createElement('a');
  a.href = `http://${DOMAIN}/${token}`;
  a.text = `${title} (created by ${username})`;
  li.appendChild(a);
  return li;
};

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

const createChat = async (api) => {
  const room = roomNameInput.value.trim();
  if (room.length < 3 || room.length > 64) return;
  const token = localStorage.getItem('token');
  try {
    await api.rooms.create({ room, token });
    roomNameInput.value = '';
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
  const websocket = await new WebSocketTransport(`ws://${DOMAIN}:8080`);
  const api = scaffold(structure, websocket);
  await setupUsername();
  await setupToken(api);
  const rooms = await api.rooms.get(); 
  for (const { title, token, username } of rooms) {
    const node = buildRoomLink(title, token, username);
    roomList.appendChild(node);
  }
  const usernameChanged = onUsernameChanged.bind(null, api);
  usernameInput.addEventListener('input', debounce(3000, usernameChanged));
  saveChatButton.addEventListener('click', createChat.bind(null, api));
  for await (const message of websocket) {
    const { title, username, token: roomsToken } = message;
    const node = buildRoomLink(title, roomsToken, username);
    roomList.appendChild(node);
  }
};

main();
