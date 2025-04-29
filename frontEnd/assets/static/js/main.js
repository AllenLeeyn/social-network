import { submitSignUp, submitLogIn, submitLogOut } from "./authentication.js";
import { addFeedbackListeners } from "./feedback.js";
import { openWebSocket } from "./messenger.js";
import { insertNewPostForm } from "./newPost.js";
import { addViewPostLinksListeners, insertPostCard } from "./post.js";
import { profileLinkHandler, addViewProfileLinksListeners } from "./profile.js";
import { templateCategoriesList, templateNoFound } from "./template.js";

document.addEventListener('DOMContentLoaded', start());
document.getElementById('logo-text').onclick = start;
document.getElementById('toLogIn').onclick = toggleView;
document.getElementById('toSignUp').onclick = toggleView;
document.getElementById('signup-btn').onclick = submitSignUp;
document.getElementById('login-btn').onclick = submitLogIn;
document.getElementById('logout-btn').onclick = submitLogOut;
document.getElementById('new-post').onclick = insertNewPostForm;

const VALIDATION_VIEW = document.getElementById("validationView");
const SIGNUP_VIEW = document.getElementById("signUpFormContainer");
const LOGIN_VIEW = document.getElementById("logInFormContainer");
const MAIN_VIEW = document.getElementById("mainView");

const PROFILE_BTN = document.getElementById("profile-btn");
PROFILE_BTN.onclick = profileLinkHandler;

const CATEGORIES_LIST = document.getElementById("categoriesList");

const FEED_DISPLAY = document.getElementById("feedDisplay");
export const POST_DISPLAY = document.getElementById("postDisplay");
export const NEW_POST_DISPLAY = document.getElementById("newPostDisplay");
export const PROFILE_DISPLAY = document.getElementById("profileDisplay");
export const MESSENGER_DISPLAY = document.getElementById("messengerDisplay");

const FEED_TAB = document.getElementById('tab-feed');
const POST_TAB = document.getElementById('tab-post');
const NEW_POST_TAB = document.getElementById('tab-new-post');
const PROFILE_TAB = document.getElementById('tab-profile');
const MESSENGER_TAB = document.getElementById('tab-messenger');

export const currentState = {
  isValid: false,
  user: "",
  id: 0,
  chat: "",
  chatID: 0,
  categories: [],
  view: SIGNUP_VIEW,
  display: FEED_DISPLAY,
  tab: FEED_TAB,
}

/*------ start ------*/
export function insertFeed(posts){
  FEED_DISPLAY.innerHTML = '';
  const container = document.createElement('tbody');
  FEED_DISPLAY.append(container);

  if (!Array.isArray(posts) || posts.length === 0){
    container.innerHTML = templateNoFound("post");
  } else {
    posts.forEach(post =>{
      insertPostCard(post, container);
    });
  }
}

export function start(){
  handleGetFetch('/posts', async (response)=>{
    if (response.ok){
      openWebSocket();
      currentState.isValid = true;
      currentState.view = MAIN_VIEW;

      const data = await response.json();
      insertFeed(data.posts);
      currentState.display = FEED_DISPLAY;
      currentState.tab = FEED_TAB;

      currentState.categories = data.categories;
      CATEGORIES_LIST.innerHTML = templateCategoriesList(currentState.categories);
      addCategoriesListeners();

      currentState.user = data.userName;
      currentState.id = data.userID;
      PROFILE_BTN.textContent = data.userName;
      PROFILE_BTN.setAttribute('href', `/profile?id=${data.userID}`)
    } else currentState.view = LOGIN_VIEW;
    renderView();
  });
}

function addCategoriesListeners(){
  const categoryItems = Array.from(CATEGORIES_LIST.children);
  categoryItems.forEach(item => addCategoriesListener(item));
}

function addCategoriesListener(item){
  const categoryHref = item.querySelector('a').getAttribute('href');
  item.onclick = (event) => {
    event.preventDefault();

    handleGetFetch(categoryHref, async response => {
      if (response.ok){
        document.getElementsByClassName('active')[0].classList.remove('active');
        item.querySelector('a').classList.add('active');
        const data = await response.json();
        insertFeed(data.posts);
        currentState.display = FEED_DISPLAY;
        currentState.tab = FEED_TAB;
        renderView();
      } else showMessage("Something went wrong. Please log in and try again.");
    });
  };
}

/*------ view functions ------*/
function toggleView(event){
  const clickedId = event.target.id;

  if (clickedId === 'toLogIn') {
    currentState.view = LOGIN_VIEW;

  } else if (clickedId === 'toSignUp') {
    currentState.view = SIGNUP_VIEW;
  }
  renderView();
}

export function renderView(){
  VALIDATION_VIEW.style.display = 'none';
  SIGNUP_VIEW.style.display = 'none';
  LOGIN_VIEW.style.display = 'none';
  MAIN_VIEW.style.display = 'none';

  if ( currentState.view === LOGIN_VIEW) {
    VALIDATION_VIEW.style.display = 'block';
    LOGIN_VIEW.style.display = 'block';

  } else if ( currentState.view === SIGNUP_VIEW) {
    VALIDATION_VIEW.style.display = 'block';
    SIGNUP_VIEW.style.display = 'block';

  } else if ( currentState.view === MAIN_VIEW) {
    MAIN_VIEW.style.display = 'flex';
    renderDisplay();
  }
}

export function renderDisplay(){
  FEED_DISPLAY.style.display = 'none';
  POST_DISPLAY.style.display = 'none';
  NEW_POST_DISPLAY.style.display = 'none';
  PROFILE_DISPLAY.style.display = 'none';
  MESSENGER_DISPLAY.style.display = 'none';

  FEED_TAB.className = '';
  POST_TAB.className = '';
  NEW_POST_TAB.className = '';
  PROFILE_TAB.className = '';
  MESSENGER_TAB.className = '';

  currentState.display.style.display = '';
  currentState.tab.className = 'selected';

  addFeedbackListeners();
  addViewPostLinksListeners();
  addViewProfileLinksListeners();
}

FEED_TAB.onclick = (event) => showDisplay(FEED_DISPLAY, event);
POST_TAB.onclick = (event) => showDisplay(POST_DISPLAY, event);
NEW_POST_TAB.onclick = (event) => showDisplay(NEW_POST_DISPLAY, event);
PROFILE_TAB.onclick = (event) => showDisplay(PROFILE_DISPLAY, event);
MESSENGER_TAB.onclick = (event) => showDisplay(MESSENGER_DISPLAY, event);

function showDisplay(display, event){
  if (currentState.display === display && display !== FEED_DISPLAY) {
    hideTab();
  } else {
    currentState.display = display;
    currentState.tab = event.target;
  }
  renderDisplay();
}

export function showTab(type, title){
  if (type === 'newPost') {
    currentState.tab = NEW_POST_TAB;
    title = "New Post";
  };
  if (type === 'post') {
    currentState.tab = POST_TAB;
    title = "Post: " + title.slice(0, 10) + "...";
  };
  if (type === 'profile') {
    currentState.tab = PROFILE_TAB;
    title = "Profile: " + title;
  };
  if (type === 'chat') {
    currentState.tab = MESSENGER_TAB;
    title = "Chat: " + title;
  };

  currentState.tab.style.display = '';
  currentState.tab.textContent = title;
}

export function hideTab(){
  if (currentState.display === MESSENGER_DISPLAY) {
    currentState.chat = "";
    currentState.chatID = 0;
  }

  currentState.display.innerHTML = '';
  currentState.tab.style.display = 'none';
  currentState.display = FEED_DISPLAY;
  currentState.tab = FEED_TAB;
}

/*------ toast message function ------*/
export function showMessage(message) {
  if (message === "") return;
  const messageDiv = document.createElement('div');
  messageDiv.classList.add('toast-message');
  messageDiv.textContent = message;
  document.body.appendChild(messageDiv);
  setTimeout(() => messageDiv.remove(), 3000);
}

/*------ handle typical JSON fetch ------*/
export function handleGetFetch(path, handler){
  fetch(path)
  .then(handler)
  .catch(error =>{
    console.error("Error:", error);
    showMessage("An error occurred. Please check your connection.");
  });
}

export function handlePostFetch(path, jsonData, message, drawFn){
  fetch(path, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    credentials: 'include',
    body: JSON.stringify(jsonData),
  })
  .then(response => {
    if (response.ok){
      showMessage(message);
      drawFn();
    } else{
      return response.json().then(errorData => {
        showMessage(errorData.message);
      });
    }
  })
  .catch(error =>{
    console.error("Error:", error);
    showMessage("An error occurred. Please check your connection.");
  });
}
