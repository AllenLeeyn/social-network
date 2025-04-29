import { NEW_POST_DISPLAY, currentState, handlePostFetch, renderDisplay, start, showTab, hideTab } from "./main.js";
import { templateNewPost } from "./template.js";

/*------ new post display ------*/
export function insertNewPostForm(event){
  event.preventDefault();
  showTab("newPost");
  NEW_POST_DISPLAY.innerHTML = '';
  const newPostElement = document.createElement('div');
  newPostElement.className = "newPost";
  newPostElement.innerHTML = templateNewPost(currentState.categories);
  NEW_POST_DISPLAY.appendChild(newPostElement);
  currentState.display = NEW_POST_DISPLAY;
  renderDisplay();
  document.getElementById('newPostForm').onsubmit = submitNewPost;
};
  
function submitNewPost(event) {
    event.preventDefault();
    const form = event.target;
    const formData = new FormData(form);
    
    const title = formData.get('title');
    const content = formData.get('content');
    const categories = formData.getAll('categories').map(c => parseInt(c));

    
    handlePostFetch('/create-post', {
      title: title,
      content: content,
      categories: categories,
    }, "Post created!", ()=>{
        hideTab("newPost");
        start()
    });
};
  