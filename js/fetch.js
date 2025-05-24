


function fetchPosts(isLoggedIn) {
    if (!isLoggedIn) {
        document.querySelector(".container").textContent =  `<h1>Welcome to the Forummmmmm</h1>`
        return;
    }
    // updateConnectedUsers(nusers);
    fetch('/api/posts')
        .then(response => response.json())
        .then(data => {
            const postsContainer = document.getElementById('posts');
            // jib dom dyal connected users and handle it 

            const createPostButton =`
                <div class="create-post">
                    <a href="/new-post" onclick="navevent(event)" class="create-post-button">Create New Post</a>
                </div>
            `;
            if (data.posts){
            postsContainer.innerHTML = `
                ${createPostButton}
                ${data.posts.map(post => `
                    <div class="post">
                        <h2><a href="/post?id=${post.ID}" onclick="navevent(event)">${post.Title}</a></h2>
                        <p>Posted by <a href="/profile?user=${post.Username}" onclick="navevent(event)">${post.Username}</a> on ${new Date(post.CreatedAt).toLocaleString()}</p>
                        <p>Comments: ${post.CommentCount}</p>
                        <p>categorie: ${post.Categories}</p>
                    </div>
                `).join('')}
            `;}else {
                postsContainer.innerHTML = `${createPostButton}`
            }
        })
        .catch(err => console.error("Error fetching posts:", err));
        convupdate(nusers)
}

function fetchcategoposts(categoID){
    console.log(categoID);
    fetch(`/api/category-posts?id=${categoID}`)
        .then(response => response.json())
        .then(data => {
            const postsContainer = document.getElementById('posts');
            const categotit = document.getElementById('catego');
            categotit.textContent = data.catego;
            convupdate(nusers);
            postsContainer.innerHTML = `
                ${data.posts.map(post => `
                    <div class="post">
                        <h2><a href="/post?id=${post.ID}" onclick="navevent(event)">${post.Title}</a></h2>
                        <p>Posted by <a href="/profile?user=${post.Username}" onclick="navevent(event)">${post.Username}</a> on ${new Date(post.CreatedAt).toLocaleString()}</p>
                        <p>Comments: ${post.CommentCount}</p>
                    </div>
                `).join('')}
            `;
        })
        .catch(err => console.error("Error fetching posts:", err));
}


function fetchProfile(user) {
    fetch(`/api/profile?user=${user}`)
        .then(response => {
            if (!response.ok) {
                throw new Error("You are notfff logged in");
            }
            return response.json();
        })
        .then(data => {
            // console.log("helelellelelele",data.isLoggedIn,data.followed.length,data);
            // console.log("lool")
            const profileContainer = document.getElementById('profile');
            var myname = localStorage.getItem('username');
            if (data.isLoggedIn) {
                console.log("logged",data)
                
                profileContainer.innerHTML = `
                    <div class="profile-info">
                        <h1>Profile Header</h1>
                        <h2>Welcome, ${data.fname} ${data.lname}!</h2>
                        ${user != myname ? 
                            (data.unfollow == true || data.status == "public") ? 
                              (data.unfollow == true ?
                                `<p>
                                  <a href="/profile?followers=${user}" onclick="navevent(event)">followers</a> : ${data.followers ? data.followers.length : 'no followers'} 
                                  <a href="/profile?followed=${user}" onclick="navevent(event)">followed</a>: ${data.followed ? data.followed.length : 'no followed'}
                                </p><a href="/profile?follow=${user}" class="button" onclick="navevent(event)">unfollow</a>` :
                                `<p>
                                  <a href="/profile?followers=${user}" onclick="navevent(event)">followers</a> : ${data.followers ? data.followers.length : 'no followers'}
                                  <a href="/profile?followed=${user}" onclick="navevent(event)">followed</a>: ${data.followed ? data.followed.length : 'no followed'}
                                </p><a href="/profile?follow=${user}" class="button" onclick="navevent(event)">follow</a>`)
                              : 
                              `<p>this profile is privet follow for moore infos</p><a href="/profile?follow=${user}" class="button" onclick="navevent(event)">follow</a>` : `<p>
                                  <a href="/profile?followers=${data.username}" onclick="navevent(event)">followers</a> : ${data.followers ? data.followers.length : 'no followers'} 
                                  <a href="/profile?followed=${data.username}" onclick="navevent(event)">followed</a>: ${data.followed ? data.followed.length : 'no followed'}
                                </p>`
                          }
                        <p><strong>Nickname:</strong> ${data.username}</p>
                        <br>
                        ${data.email ? `<p><strong>Email:</strong> ${data.email}</p>` : ``}
                    </div>
                    <br>
                    <div id="profile-posts">
                    <h3> Posts </h3>
                    ${data.post ? data.posts.map(post => `
                        <h2><a href="/post?id=${post.ID}" onclick="navevent(event)">${post.Title}</a></h2>
                        <br>
                    </div>
                `).join('') : ``}
                `;
                convupdate(nusers);
            } else {
                profileContainer.innerHTML = `
                    <div class="login-message">
                        <h2>You are not logged in!!!!!!</h2>
                        <a href="/login" onclick="navevent(event)" class="login-button">Go to Login</a>
                    </div>
                `;
            }
        })
        .catch(error => {
            const profileContainer = document.getElementById('profile');
            profileContainer.innerHTML = `
                <div class="login-message">
                    <h2>You are not logged inin!</h2>
                    <a href="/login" onclick="navevent(event)" class="login-button">Go to Login</a>
                </div>
            `;
        });
}


function fetchPost(postID) {
    fetch(`/api/post?id=${postID}`)
        .then(response => response.json())
        .then(data => {
            const postContainer = document.getElementById('post');
            console.log(data.post.Attachement)
            postContainer.innerHTML = `
                <p><strong>Posted on:</strong> ${new Date(data.post.CreatedAt).toLocaleString()} | <strong>Categories:</strong> ${data.post.Categories} | <strong>Created By:</strong> ${data.post.Username}</p>
                <h1>${data.post.Title}</h1>
                ${data.post.Attachement ? `
                    <div class="post-attachment">
                        ${data.post.Attachement.match(/\.(jpg|jpeg|png|gif)$/i) ? 
                            `<img src="${data.post.Attachement}">` : 
                            `<a href="${data.post.Attachement}" target="_blank">Download Attachment</a>`
                        }
                    </div>
                ` : ''}
                <p>${data.post.Content}</p>
                <h3>Comments</h3>
                <ul>
                ${data.comments ? `
                    ${data.comments.map(comment => `
                        <li>
                            <strong>${comment.Username}:</strong> ${comment.Content} (Posted on: ${new Date(comment.CreatedAt).toLocaleString()})
                        </li>
                    `).join('')}
                 ` : ``}
                 </ul>
                ${data.isLoggedIn ? `
                    <form id="comment-form">
                        <input type="hidden" name="post_id" value="${data.post.ID}">
                        <textarea name="comment" placeholder="Add your comment" required></textarea>
                        <br>
                        <button type="submit">Submit</button>
                    </form>
                ` : `
                    <p><a href="/register" onclick="navevent(event)">Register</a> to add a comment.</p>
                `}

                <a href="/" onclick="navevent(event)">Back to Home</a>
            `;

            setupFormListeners();
        })
        .catch(err => console.error("Error fetching post:", err));
}

function fetchCategories() {
    // updateConnectedUsers(nusers);
    convupdate(nusers);
    fetch('/api/categories')
        .then(response => response.json())
        .then(data => {
            const container = document.getElementById('categories');
            if (!data.categories || data.categories.length === 0) {
                container.innerHTML = `<p>No categories found</p>`;
                return;
            }
            const categoriesContainer = document.getElementById('categories');
            categoriesContainer.innerHTML = `
                <div class="grid">
                    ${data.categories.map(category => `
                        <div class="category-item">
                            <a href="/category-posts?id=${category.id}" onclick="navevent(event)">
                                <h3>${category.name}</h3>
                                <p>Related Posts: ${category.postCount} Posts</p>
                            </a>
                        </div>
                    `).join('')}
                </div>
            `;
        })
        .catch(err => console.error("Error fetching categories:", err));
}