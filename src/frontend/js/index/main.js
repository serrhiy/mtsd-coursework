'use strict';

import WebSocketTransport from "../WebSocketTransport.js";
import structure from "./structure.js";
import scaffold from "../scaffold.js";

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

const main = async () => {
  const websocket = await new WebSocketTransport('ws://127.0.0.1:8080');
  const api = scaffold(structure, websocket);
  console.log(await api.users.exists("asd"));
  saveUsernameButton.addEventListener('click', onUsernameButton.bind(null, api));
  saveChatButton.addEventListener('click', onChatButton.bind(null, api));
};

main();