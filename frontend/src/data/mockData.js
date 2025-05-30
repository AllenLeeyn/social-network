// ====================
// Posts
// ====================
export const samplePosts = [
        {
        id: 1,
        title: 'Post Title 1',
        author: 'UserA',
        snippet: 'This is a snippet of the first post...'
    },
    {
        id: 2,
        title: 'Post Title 2',
        author: 'UserB',
        snippet: 'This is a snippet of the second post...',
        online: true,
        unread: true
    },
    {
        id: 3,
        title: 'Post Title 3',
        author: 'UserC',
        snippet: 'This is a snippet of the third post...'
    }
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
        timestamp: "2024-05-09T18:00:00Z"
    },
    {
        id: 2,
        postId: 2,
        author: "UserA",
        content: "Thanks for sharing!",
        timestamp: "2024-05-09T18:05:00Z"
    },
    {
        id: 3,
        postId: 1,
        author: "UserB",
        content: "Interesting thoughts.",
        timestamp: "2024-05-09T18:10:00Z"
    }
];

  // ====================
  // Categories
  // ====================
export const sampleCategories = [
    { id: 1, name: "Technology" },
    { id: 2, name: "Health" },
    { id: 3, name: "Travel" }
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
    online: true
    },
    {
    id: 2,
    username: "UserB",
    fullName: "Bob Brown",
    avatar: "/avatars/bob.png",
    online: true,
    unread: true
    },
    {
    id: 3,
    username: "UserC",
    fullName: "Charlie Clark",
    avatar: "/avatars/charlie.png"
    }
];

// ====================
// Groups
// ====================
export const sampleGroups = [
    { id: 1, name: "React Enthusiasts" },
    { id: 2, name: "Travel Buddies" },
    { id: 3, name: "Book Club" }
];

// ====================
// Connections
// ====================
export const sampleConnections = [
    { id: 1, username: "UserD", fullName: "Dana Doe" },
    { id: 2, username: "UserE", fullName: "Evan Evans" }
];

<<<<<<< HEAD

export const sampleConversations = [
  {
    id: 'conv1',
    name: 'Alice Smith',
    type: 'individual',
    unread: 2
  },
  {
    id: 'conv2',
    name: 'Bob Johnson',
    type: 'individual',
    unread: 0
  },
  {
    id: 'conv3',
    name: 'Project Team',
    type: 'group',
    unread: 0
  },
  {
    id: 'conv4',
    name: 'Friends Group',
    type: 'group',
    unread: 3
  }
];

export const sampleMessages = [
  {
    id: 'msg1',
    content: 'Hey, how are you?',
    senderId: 'user1',
    timestamp: '2025-05-20T10:00:00Z',
    conversationId: 'conv1'
  },
  {
    id: 'msg2',
    content: 'Are we meeting tomorrow?',
    senderId: 'user2',
    timestamp: '2025-05-20T11:30:00Z',
    conversationId: 'conv1'
  },
  {
    id: 'msg3',
    content: 'Check out this link!',
    senderId: 'user3',
    timestamp: '2025-05-21T09:15:00Z',
    conversationId: 'conv3'
  }
];

export const mockNotifications = [
  {
    id: 'notif-1',
    type: 'follow_request',
    isRead: false,
    timestamp: '2025-05-27T08:30:00Z',
    fromUser: {
      id: 'user-2',
      name: 'Alice Johnson',
      avatar: '/avatars/alice.jpg'
    },
    message: 'Alice Johnson wants to follow you.',
    actions: ['accept', 'decline']
  },
  {
    id: 'notif-2',
    type: 'group_invitation',
    isRead: false,
    timestamp: '2025-05-27T07:50:00Z',
    group: {
      id: 'group-1',
      name: 'React Enthusiasts',
      avatar: '/groups/react.png'
    },
    fromUser: {
      id: 'user-3',
      name: 'Bob Lee',
      avatar: '/avatars/bob.jpg'
    },
    message: 'Bob Lee invited you to join React Enthusiasts.',
    actions: ['accept', 'decline']
  },
  {
    id: 'notif-3',
    type: 'group_join_request',
    isRead: true,
    timestamp: '2025-05-26T20:15:00Z',
    group: {
      id: 'group-2',
      name: 'Next.js Masters',
      avatar: '/groups/nextjs.png'
    },
    fromUser: {
      id: 'user-4',
      name: 'Charlie Kim',
      avatar: '/avatars/charlie.jpg'
    },
    message: 'Charlie Kim requested to join your group Next.js Masters.',
    actions: ['accept', 'decline']
  },
  {
    id: 'notif-4',
    type: 'group_event',
    isRead: false,
    timestamp: '2025-05-27T06:00:00Z',
    group: {
      id: 'group-1',
      name: 'React Enthusiasts',
      avatar: '/groups/react.png'
    },
    event: {
      id: 'event-1',
      name: 'React 19 Launch Party'
    },
    fromUser: {
      id: 'user-5',
      name: 'Dana White',
      avatar: '/avatars/dana.jpg'
    },
    message: 'A new event "React 19 Launch Party" was created in React Enthusiasts.',
    actions: ['view']
  }
];


// src/data/mockGroups.js

export const mockGroups = [
  {
    id: 1,
    title: "React Enthusiasts",
    description: "A group for React lovers",
    members: ["alice", "bob"]
  },
  {
    id: 2,
    title: "Next.js Learners",
    description: "Learning Next.js together",
    members: ["carol"]
  },
  {
    id: 3,
    title: "Open Source Contributors",
    description: "Contribute to open source projects",
    members: ["dave", "eve"]
  }
];

export const mockInvitations = [
  {
    id: 101,
    groupId: 1,
    groupTitle: "React Enthusiasts",
    fromUser: "alice",
    toUser: "frank",
    status: "pending"
  },
  {
    id: 102,
    groupId: 2,
    groupTitle: "Next.js Learners",
    fromUser: "carol",
    toUser: "grace",
    status: "pending"
  }
];

export const mockEvents = [
  {
    id: 201,
    groupId: 1,
    groupTitle: "React Enthusiasts",
    title: "React Hooks Deep Dive",
    description: "An in-depth look at React Hooks",
    dateTime: "2025-06-01T18:00:00Z",
    rsvps: { alice: "going", bob: "not_going" }
  },
  {
    id: 202,
    groupId: 3,
    groupTitle: "Open Source Contributors",
    title: "Monthly OSS Meeting",
    description: "Discuss open source projects and contributions",
    dateTime: "2025-06-05T15:00:00Z",
    rsvps: { dave: "going", eve: "going" }
  }
];


// In mockGroups.js
export const mockPosts = [
  {
    id: 1,
    groupId: 1,
    title: "Welcome to React Enthusiasts!",
    content: "Let's share our favorite React tips here.",
    author: "alice",
  },
  {
    id: 2,
    groupId: 2,
    title: "Next.js Meetup",
    content: "Who wants to join a virtual meetup next week?",
    author: "carol",
  }
=======
export const sampleFollowers = [
    { id: 1, username: "UserA", fullName: "Alice Anderson", avatar: "/avatars/alice.png" },
    { id: 2, username: "UserB", fullName: "Bob Brown", avatar: "/avatars/bob.png" },
    { id: 3, username: "UserC", fullName: "Charlie Clark", avatar: "/avatars/charlie.png" },
];

export const sampleFollowing = [
    { id: 4, username: "UserD", fullName: "David Davis", avatar: "/avatars/david.png" },
    { id: 5, username: "UserE", fullName: "Emma Evans", avatar: "/avatars/emma.png" },
    { id: 6, username: "UserF", fullName: "Frank Foster", avatar: "/avatars/frank.png" },
];

export const myActivity = [
    { id: 1, name: "My Posts" },
    { id: 2, name: "My Group Posts" },
>>>>>>> gigi
];
