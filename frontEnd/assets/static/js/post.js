import { currentState, POST_DISPLAY, renderDisplay, handleGetFetch, handlePostFetch, showMessage, showTab } from "./main.js";
import { templateCommentCard, templateCommentForm, templateNoFound, templatePostCard } from "./template.js";

export function insertPostCard(post, container){
    const row = document.createElement('tr');
    const cell = document.createElement('td');

    const postElement = document.createElement('div');
    postElement.classList.add('post-card');
    postElement.innerHTML = templatePostCard(post);

    cell.appendChild(postElement);
    row.appendChild(cell);
    container.appendChild(row);

    return postElement;
}

export function addViewPostLinksListeners(){
    const postLinks = document.querySelectorAll(".post-content h3 a");
    const commentBtns = document.querySelectorAll(".fas.fa-comment");

    postLinks.forEach(link =>{
        link.onclick = postLinkHandler;
    });
    commentBtns.forEach(link =>{
        link.onclick = postLinkHandler;
    });
}

function postLinkHandler(event){
    event.preventDefault();
    const path = event.target.getAttribute('href');
    
    handleGetFetch(path, insertPostWithComments);
}

function getPost(data){
    const container = document.createElement('tbody');
    const postElement = insertPostCard(data.post, container);
    postElement.innerHTML += templateCommentForm(data.post);

    if (!Array.isArray(data.comments) || data.comments.length === 0){
        postElement.innerHTML += templateNoFound("comment");
    } else {
        data.comments.forEach(comment => {
            postElement.innerHTML += templateCommentCard(comment);
        });
    }
    return container;
}

function submitComment(event){
    const postID = parseInt(event.target.getAttribute('data-id'));
    const commentInput = document.getElementById('comment-input')
    const commentText = commentInput.value.trim()

    if (!commentText) {
        showMessage('Comment cannot be empty!')
        return;
    }
    handlePostFetch(`/create-comment`, { 
        postID: postID,
        content: commentText,
    }, "Comment created!", ()=>{
        handleGetFetch(`/post?id=${postID}`, insertPostWithComments)
    });
};

async function insertPostWithComments(response){
    if (response.ok) {
        POST_DISPLAY.innerHTML = "";
        const data = await response.json();
        POST_DISPLAY.appendChild(getPost(data));
        currentState.display = POST_DISPLAY;
        showTab("post", data.post.title);
        renderDisplay();
        document.getElementById('submit-comment').onclick = submitComment;
    } else showMessage("Something went wrong. Please log in and try again.");
}
