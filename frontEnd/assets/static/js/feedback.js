import { handlePostFetch } from "./main.js";

export function addFeedbackListeners(){
  const likeButtons = document.querySelectorAll(".like-button");
  likeButtons.forEach((button) => {
    button.onclick = () => feedbackHandler("like", button);
  });
  
  const dislikeButtons = document.querySelectorAll(".dislike-button");
  dislikeButtons.forEach((button) => {
    button.onclick = () => feedbackHandler("dislike", button);
  });
}

// Unified function to handle both like and dislike actions
function feedbackHandler(action, button) {
  const forType = button.getAttribute("data-for");
  const parentID = parseInt(button.getAttribute("data-id"), 10);
  const parentElement = button.closest('.post-actions')

  const currentState = parseInt(parentElement.getAttribute('data-state'), 10);
  let newState = currentState;

  // Get like and dislike button elements for the post
  const likeButton = parentElement.children[0];;
  const likeCountSpan = likeButton.querySelector("span");
  let newLikeCount = parseInt(likeCountSpan.textContent);

  const dislikeButton = parentElement.children[1];;
  const dislikeCountSpan = dislikeButton.querySelector("span");
  let newDislikeCount = parseInt(dislikeCountSpan.textContent);

  [newState, newLikeCount, newDislikeCount] = toggleFeedback(action, currentState, newLikeCount, newDislikeCount);

  handlePostFetch(`/feedback`, {
      tgt: forType,
      parentID: parentID,
      rating: newState,
    }, "", ()=>{
      parentElement.setAttribute('data-state', newState);
      likeCountSpan.textContent = newLikeCount;
      dislikeCountSpan.textContent = newDislikeCount;
    });
}

function toggleFeedback(action, currentState, newLikeCount, newDislikeCount){
  let toggler = 1;
  let counter = newLikeCount;
  let altCounter = newDislikeCount;
  let newState = 0;

  if (action === "dislike"){
    toggler = -1;
    counter = newDislikeCount;
    altCounter = newLikeCount;
  }

  // remove current feedback if action matches current feedback
  if (currentState === toggler) {
    counter = counter - 1;
    newState = 0;
  } else {
    if (currentState === -toggler) {
      altCounter = altCounter - 1;
    }
    counter = counter + 1;
    newState = toggler;
  }
  newLikeCount = (action === "like") ? counter : altCounter;
  newDislikeCount = (action === "like") ? altCounter : counter;

  return [newState, newLikeCount, newDislikeCount];
}