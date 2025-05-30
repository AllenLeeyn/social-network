// ====================
// Posts
// ====================
export const samplePosts = [
  {
    id: 1,
    title: "Post Title 1",
    author: "UserA",
    snippet: "This is a snippet of the first post...",
  },
  {
    id: 2,
    title: "Post Title 2",
    author: "UserB",
    snippet: "This is a snippet of the second post...",
    online: true,
    unread: true,
  },
  {
    id: 3,
    title: "Post Title 3",
    author: "UserC",
    snippet: "This is a snippet of the third post...",
  },
];

// ====================
// Comments
// ====================
export const sampleComments = [
  {
    id: 1,
    postId: 2,
    author: "UserC",
    content: "Great post!",
    timestamp: "2024-05-09T18:00:00Z",
  },
  {
    id: 2,
    postId: 2,
    author: "UserA",
    content: "Thanks for sharing!",
    timestamp: "2024-05-09T18:05:00Z",
  },
  {
    id: 3,
    postId: 1,
    author: "UserB",
    content: "Interesting thoughts.",
    timestamp: "2024-05-09T18:10:00Z",
  },
];

// ====================
// Categories
// ====================
export const sampleCategories = [
  { id: 1, name: "Technology" },
  { id: 2, name: "Health" },
  { id: 3, name: "Travel" },
];

// ====================
// Users
// ====================
export const sampleUsers = [
  {
    id: 1,
    username: "UserA",
    fullName: "Alice Anderson",
    avatar: "/avatars/alice.png",
    online: true,
  },
  {
    id: 2,
    username: "UserB",
    fullName: "Bob Brown",
    avatar: "/avatars/bob.png",
    online: true,
    unread: true,
  },
  {
    id: 3,
    username: "UserC",
    fullName: "Charlie Clark",
    avatar: "/avatars/charlie.png",
  },
];

// ====================
// Groups
// ====================
export const sampleGroups = [
  { id: 1, name: "React Enthusiasts" },
  { id: 2, name: "Travel Buddies" },
  { id: 3, name: "Book Club" },
];

// ====================
// Connections
// ====================
export const sampleConnections = [
  { id: 1, username: "UserD", fullName: "Dana Doe" },
  { id: 2, username: "UserE", fullName: "Evan Evans" },
];

export const sampleFollowers = [
  {
    id: 1,
    username: "UserA",
    fullName: "Alice Anderson",
    avatar: "/avatars/alice.png",
  },
  {
    id: 2,
    username: "UserB",
    fullName: "Bob Brown",
    avatar: "/avatars/bob.png",
  },
  {
    id: 3,
    username: "UserC",
    fullName: "Charlie Clark",
    avatar: "/avatars/charlie.png",
  },
];

export const sampleFollowing = [
  {
    id: 4,
    username: "UserD",
    fullName: "David Davis",
    avatar: "/avatars/david.png",
  },
  {
    id: 5,
    username: "UserE",
    fullName: "Emma Evans",
    avatar: "/avatars/emma.png",
  },
  {
    id: 6,
    username: "UserF",
    fullName: "Frank Foster",
    avatar: "/avatars/frank.png",
  },
];

export const myActivity = [
  { id: 1, name: "My Posts" },
  { id: 2, name: "My Group Posts" },
];
