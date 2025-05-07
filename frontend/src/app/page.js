// This is our homepage
import React from 'react';

export default function HomePage() {
const samplePosts = [
{ id: 1, title: 'Post Title 1', author: 'UserA', snippet: 'This is a snippet of the first post...' },
{ id: 2, title: 'Post Title 2', author: 'UserB', snippet: 'This is a snippet of the second post...' },
{ id: 3, title: 'Post Title 3', author: 'UserC', snippet: 'This is a snippet of the third post...' },
];

return (
    <main>
        <h1>Welcome to grit:Hub!</h1>
        <p>Your place to connect and share.</p>
        <div style={{ marginBottom: '1rem' }}>
        <button style={{ marginRight: '1rem' }}>Sign Up</button>
        <button>Log In</button>
        </div>
        <section>
        <h2>Latest Posts</h2>
        <div>
            {samplePosts.map(post => (
            <div key={post.id} className="post-card">
                <h3>{post.title}</h3>
                <p><em>by {post.author}</em></p>
                <p>{post.snippet}</p>
            </div>
            ))}
        </div>
        </section>
    </main>
    );
}






/* <main>
    <h1>Welcome to SocialNet!</h1>
    <p>Your place to connect and share.</p>
    <div>
        [Sign Up Button]    [Log In Button]
    </div>
    <section>
        <h2> Latest Posts</h2>
        <div>
            [PostCard]  [PostCard]  [PostCard]
        </div>
    </section>
</main> */
