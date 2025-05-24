const app = document.getElementById('app');
const navbar = document.getElementById('navbar');
const authLink = document.getElementById('authLink');
let isProgrammaticNavigation = false;
let isLoggedIn = false; 

const routes = {
    '/': `
    <div class="main-container">
        <div class="chat-sidebar">
            <h2>Conversations</h2>
            <div id="conversations-list"></div>
        </div>
        <div class="container">
            <div id="posts"></div>
        </div>
    <div class="chat-sidebar">
    `,
    '/register': `
<div class="containerRegister">
    <div class="form-container">
        <h1>Register</h1>
        <form id="registerForm">
            <div class="form-row">
                <div class="form-group">
                    <label for="username">Nickname:</label>
                    <input type="text" name="username" id="username" required>
                </div>
                <div class="form-group">
                    <label for="email">Email:</label>
                    <input type="email" name="email" id="email" required>
                </div>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label for="first">First Name:</label>
                    <input type="text" name="first" id="first" required>
                </div>
                <div class="form-group">
                    <label for="last">Last Name:</label>
                    <input type="text" name="last" id="last" required>
                </div>
            </div>

            <div class="form-row">
                <div class="form-group">
                    <label for="age">Your Age:</label>
                    <input type="number" id="age" name="age" min="0" max="150" required>
                </div>
                <div class="form-group">
                    <label>Gender:</label>
                    <div class="gender-options">
                        <label><input type="radio" name="gender" value="male" required> Male</label>
                        <label><input type="radio" name="gender" value="female"> Female</label>
                    </div>
                </div>
            </div>
            
            <div class="form-row">
                <div class="form-group">
                    <label for="birthdate">Birthday:</label>
                    <input type="date" name="birthdate" id="birthdate" required>
                </div>
                <div class="form-group">
                    <label for="avatar">Avatar:</label>
                    <input type="file" name="avatar" id="avatar" accept="image/*">
                </div>
            </div>

            <div class="form-group full">
                <label for="about">About Me:</label>
                <textarea name="about" id="about" rows="4" placeholder="Tell us something about yourself..."></textarea>
            </div>

            <div class="form-group full">
                <label for="password">Password:</label>
                <input type="password" name="password" id="password" required>
            </div>

            <button type="submit">Register</button>
        </form>
        <a href="/login" onclick="navevent(event)">Already have an account? Login here.</a>
    </div>
</div>

    `,
    '/login': `
        <div class="containerLogin">
            <div class="form-container">
                <h1>Login</h1>
                <form id="loginForm">
                    <label for="username">Nickname oder Email:</label>
                    <input type="text" name="username" id="username" required>
                    <br>
                    <label for="password">Password:</label>
                    <input type="password" name="password" id="password" required>
                    <br>
                    <button type="submit">Login</button>
                </form>
                <a href="/register" onclick="navevent(event)">Don't have an account? Register here.</a>
            </div>
        </div>
    `,
    '/new-post': `
    <div class="main-container">
        <div class="chat-sidebar">
            <h2>Conversations</h2>
            <div id="conversations-list"></div>
        </div>
    <div class="container">
        <div class="form-container">
            <h1>Create New Post</h1>
            <form id="addPostForm">
    <label for="title">Title:</label>
    <input type="text" name="title" id="title" required>
    <br>

    <label for="content">Content:</label>
    <textarea name="content" id="content" required></textarea>
    <br>

    <label for="category">Associate Categories to Post:</label>
    <div class="categories" id="category">
        <label><input type="checkbox" name="categories[]" value="Technology"> Technology</label>
        <label><input type="checkbox" name="categories[]" value="Health"> Health</label>
        <label><input type="checkbox" name="categories[]" value="Travel"> Travel</label>
        <label><input type="checkbox" name="categories[]" value="Sports"> Sports</label>
        <label><input type="checkbox" name="categories[]" value="Gaming"> Gaming</label>
        <label><input type="checkbox" name="categories[]" value="Food"> Food</label>
        <label><input type="checkbox" name="categories[]" value="Science"> Science</label>
        <label><input type="checkbox" name="categories[]" value="Fashion"> Fashion</label>
    </div>
    <br>

    <!-- New: Private Post Checkbox -->
        <label>
            <input type="checkbox" id="is_private" name="is_private">
            Make this post private
        </label>
        <br>

        <!-- New: File Attachment -->
        <label for="attachment">Attach a file:</label>
        <input type="file" name="attachment" id="attachment" accept=".jpg,.jpeg,.png,.pdf,.doc,.docx,.txt">
        <br>

        <button type="submit">Submit</button>
    </form>

        </div>
    </div>
    </div>
`,
    '/profile': `
    <div class="main-container">
        <div class="chat-sidebar">
            <h2>Conversations</h2>
            <div id="conversations-list"></div>
        </div>
        <div class="container">
            <div id="profile"></div>
        </div>
    </div>
`,
'/categories': `
    <div class="main-container">
        <div class="chat-sidebar">
            <h2>Conversations</h2>
            <div id="conversations-list"></div>
        </div>
        <div class="container">
            <div id="categories"></div>
        </div>
    </div>
    `,
'/post': `
    <div class="main-container">
        <div class="chat-sidebar">
            <h2>Conversations</h2>
            <div id="conversations-list"></div>
        </div>
        <div class="container">
            <div id="post"></div>
        </div>
    </div>
    `,
'/chat': `
<div class="main-container">

        <div class="chat-sidebar">
            <h2>Conversations</h2>
            <div id="conversations-list"></div>
        </div>
    <div class="chat-container">
        <div class="chat-main">
            <div class="chat-header">
                <h3 id="chat-with-user">Select a conversation</h3>
            </div>
            <div id="chat-messages" class="chat-messages"></div>
            <div class="chat-input">
                <input type="text" id="chatMessage" placeholder="Type your message..." disabled>
                <button id="sendMessage" disabled>Send</button>
            </div>
        </div>
    </div>
    </div>
`,
'/category-posts':`
<div class="main-container">
 <div class="chat-sidebar">
            <h2>Conversations</h2>
            <div id="conversations-list"></div>
        </div>
        <div class="container">
            <h1 id="catego"></h1>
            <hr />
            <div id="posts"></div>
        </div>
    </div>
    `,
};
let previousPath = '';

function navevent(event) {
    event.preventDefault();
    const link = event.target.closest('a');
    if (!link) return;
    
    // console.log(link.getAttribute('href'));/
    navigate(link.getAttribute('href'));
}

function navigate(path ,{ back = true } = {}) {
    isProgrammaticNavigation = true;
    var basePath = ''
    console.log(path,"hifirstyytt")
    basePath = path.split('?')[0];
    if(routes[basePath]){
        app.innerHTML = routes[basePath];
        if (back){
            history.pushState({ path: path }, path, path);
        }
    const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true';
    const queryParams = new URLSearchParams(window.location.search);
    console.log(path, basePath,"hi");
    if (!isLoggedIn && basePath != "/login" && basePath != "/register") {
        app.innerHTML = routes["/login"];
    }
    updateAuthLink();
    setupFormListeners();
    
    if ((basePath === '/' && isLoggedIn) || ((basePath == "/login" || basePath == "/register") && isLoggedIn)) {
        app.innerHTML = routes["/"]
        fetchPosts(isLoggedIn);
    } else if (basePath === '/profile' && isLoggedIn) {
        const profileuser = queryParams.get('user');
        const profilefollow = queryParams.get('follow')
        const pfollowers = queryParams.get('followers')
        const pfollowed = queryParams.get('followed')
        if (profileuser) {
            fetchProfile(profileuser);
        }else if (profilefollow) {
            follower(1,profilefollow);
        }else if (pfollowers){
            follower(0,pfollowers)
        }else if (pfollowed){
            follower(-1,pfollowed); 
        }
    } else if (basePath === '/categories' && isLoggedIn) {
        fetchCategories();
    } else if (basePath === '/post' && isLoggedIn) {
        const postID = queryParams.get('id');
        if (postID){
        fetchPost(postID);
        convupdate(nusers);
        }else {
            navigate("/");
        }

    } else if (basePath === '/category-posts' && isLoggedIn) {
        const categoID = queryParams.get('id');
        if (categoID){
        fetchcategoposts(categoID);
        }else {
            navigate("/");
        }
    } else if (basePath === '/chat' && isLoggedIn) {
        convupdate(nusers);
        if (currentChatUser) {
            document.getElementById('chatMessage').disabled = false;
            document.getElementById('sendMessage').disabled = false;
        }
    }else if  (basePath == "/new-post" ) {
       convupdate(nusers)
    }
}else {
    app.innerHTML = `<h1>404 Not Found</h1>`
}
    isProgrammaticNavigation = false;
}

function follower(number,f) {
    // fetchProfile(f)
    console.log()
    if(number == 0 || number == -1) {
        fetch(`/api/profile?user=${f}`)
        .then(response => {
            if (!response.ok) {
                throw new Error("You are not logged innn");
            }
            return response.json();
        })
        .then(data => {
            console.log("hiii its me")
            const cont = getElementByid('profile-posts');
            cont.innerHTML = ``
            if (number == 0) {
                cont.innerHTML =`
                <ul>
                ${data.followers.map(fl => `<li><a href="/profile?user=${fl}" onclick="navevent(event)">${fl}</a></li>`).join('')}
              </ul>
            `;
            }else {
                cont.innerHTML =`
                <ul>
                ${data.followed.map(fd => `<li><a href="/profile?user=${fd}" onclick="navevent(event)">${fd}</a></li>`).join('')}
              </ul>
            `;
            }
    }).catch(error => {
        const profileContainer = document.getElementById('profile');
        profileContainer.innerHTML = `
            <div class="login-message">
                <h2>You are not logged innnn!</h2>
                <a href="/login" onclick="navevent(event)" class="login-button">Go to Login</a>
            </div>
        `;
    });
}else {
    fetch(`/api/follow?id=${f}`)
    .then(response => {
        if (!response.ok) {
            throw new Error("You are not logged in");
        }
        return response.json();
    })
    .then(data => { 
        alert(data.status);

    }).catch(error => {
        const profileContainer = document.getElementById('profile');
        profileContainer.innerHTML = `
            <div class="login-message">
                <h2>You are not logged inff!</h2>
                <a href="/login" onclick="navevent(event)" class="login-button">Go to Login</a>
            </div>
        `;
    }); 
}
}


document.body.addEventListener('click', e => {
    if (e.target.matches('[data-link]')) {
        e.preventDefault();
        const path = e.target.getAttribute('href');
        if (path) {
            console.log("booooooofffyyyyyyyyyyy");
            navigate(path);
        } else {
            console.error("Link clicked with no href!");
            navigate('/errors?code=400');
        }
    }
});

window.addEventListener('popstate', e => {
    if (!isProgrammaticNavigation && e.state && e.state.path) { // Double check e.state
        navigate(e.state.path, { pushHistory: false });
    } else {
        navigate('/', { pushHistory: false }); // Safe default
    }
});

function updateAuthLink() {
    const isLoggedIn = localStorage.getItem('isLoggedIn') === 'true';

    if (isLoggedIn) {
        const usernm = localStorage.getItem('username');
        navbar.innerHTML = `
            <a href="/" onclick="navevent(event)">Home</a>
            <a href="/categories" onclick="navevent(event)">categories</a>
            <a href="/profile?user=${usernm}" onclick="navevent(event)">Profile</a>
            <a href="#" onclick="handleLogout()">Logout</a>
        `;
    }else {
        navbar.innerHTML = `
            <a href="/login" onclick="navevent(event)">Login</a>
            <a href="/register" onclick="navevent(event)">Register</a>
        `;
    }
}


function handleLogout() {
    localStorage.removeItem('isLoggedIn');
    fetch('http://localhost:8080/api/logout', { method: 'POST' })
    .then(()=>{
        if (ws) {
            ws.onclose();
        }
        navigate("/login");
    })
    .catch(err => console.error("Logout error:", err));
}

async function checkAuthStatus() {
    try {
        const response = await fetch('/api/auth');
        if (response.ok) {
            const data = await response.json();
            if (data.isLoggedIn) {
                localStorage.setItem('isLoggedIn', 'true');
                if (data.username) {
                    localStorage.setItem('username', data.username);
                }
                connectWebSocket();
                return true;
            }
        }
    } catch (err) {
        console.error("Auth check error:", err);
    }
    localStorage.removeItem('isLoggedIn');
    localStorage.removeItem('username');
    return false;
}

document.addEventListener('DOMContentLoaded', async () => {
    await checkAuthStatus();
    const path = window.location.pathname + window.location.search; //hnaaaaaaaaaaaaaaaa
        navigate(path);
});


function setupFormListeners() {
    const registerForm = document.getElementById('registerForm');

    if (registerForm) {
        registerForm.addEventListener('submit', function(event) {
            event.preventDefault();
    
            const formData = new FormData(registerForm);
            const avatarFile = formData.get('avatar');
            const avatar = avatarFile ? avatarFile.name : '';
            const data = {
                username: formData.get('username'),
                email: formData.get('email'),
                password: formData.get('password'),
                age: formData.get('age'),
                gender: formData.get('gender'),
                fname: formData.get('first'),
                lname: formData.get('last'),
                birthdate: formData.get('birthdate'),
                avatar: avatar,
                aboutme: formData.get('about'),
                status: "public"
            };
            
            fetch('http://localhost:8080/api/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            })
            .then(response => {
                if (!response.ok) {
                    return response.json().then(err => { throw new Error(err.error || "Registration failed"); });
                }
                return response.json();
            })
            .then(data => {
                console.log("Registration successful:", data);
                alert("Registration successful! Redirecting to login page.");
                window.location.href = "/login"; 
            })
            .catch(err => {
                console.error("Registration error:", err);
                alert(err.message); 
            });
        });
    }
    

    const loginForm = document.getElementById('loginForm');

    if (loginForm) {
        loginForm.addEventListener('submit', function(event) {
            event.preventDefault();
            
            const formData = new FormData(loginForm);
            const data = {
                username: formData.get('username'),
                password: formData.get('password')
            };
    
            fetch('http://localhost:8080/api/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            })
            .then(response => {
                if (!response.ok) {
                    return response.json().then(err => { throw new Error(err.error || "Login failed"); });
                }
                return response.json();
            })
            .then(data => {
                localStorage.setItem('isLoggedIn', 'true');
                localStorage.setItem('username', data.username);//removie
                navigate('/');
                connectWebSocket();
            })
            .catch(err => {
                console.error("Login error:", err);
                alert(err.message); 
            });
        });
    }
    

    const addPostForm = document.getElementById('addPostForm');

    if (addPostForm) {
        addPostForm.addEventListener('submit', async function(event) {
            event.preventDefault();
        
            const formData = new FormData(addPostForm);
            const postFile = formData.get('attachment');
            const attachmentfile = postFile ? postFile.name : '';
            const title = formData.get('title');
            const content = formData.get('content');
            const categories = formData.getAll('categories[]');
            const isPrivate = formData.get('is_private') === 'on';
            const statusPost = isPrivate ? 'private' : 'public';
            const attachment = attachmentfile
            let attachmentPath = null;
            const file = formData.get('attachment');
            console.log(file,"fillle")
            if (file && file.size > 0) {
                const uploadForm = new FormData();
                uploadForm.append('file', file);
        
                try {
                    const uploadResponse = await fetch('http://localhost:8080/api/upload', {
                        method: 'POST',
                        body: uploadForm
                    });
                    const uploadResult = await uploadResponse.json();
                    if (uploadResult.path) {
                        attachmentPath = uploadResult.path;
                        console.log(attachmentPath,"attttachemnt")
                    }
                } catch (err) {
                    console.error('File upload failed:', err);
                }
            }
            // console.log(attachmentPath)
            const data = {
                title,
                content,
                category: categories,
                statusPost,
                attachment,
                attachment: attachmentPath
            };
            
            fetch('http://localhost:8080/api/add-post', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            })
            .then(response => response.json())
            .then(data => {
                alert(data.message);
                window.location.href = '/';
            })
            .catch(err => console.error("Post submission error:", err));
        });
        
    }
   
    const addcommentform = document.getElementById('comment-form');
    if (addcommentform) {
        // console.log("hellllllllllo comment");
        addcommentform.addEventListener('submit', function(event) {
        event.preventDefault();
        const formData = new FormData(addcommentform);
        const postID = formData.get('post_id');
        const comment = formData.get('comment');
        fetch('/api/add-comment', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ post_id: postID, comment })
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                const commentList = document.querySelector('ul');
                const newComment = document.createElement('li');
                newComment.innerHTML = `<strong>${data.username}:</strong> ${data.content} (Posted on: ${new Date(data.createdAt).toLocaleString()})`;
                console.log("yahooo");
                commentList.appendChild(newComment);
                addcommentform.querySelector('textarea').value = '';
            } else {
                console.error('Error adding comment!');
            }
        })
        .catch(error => console.error('Error:', error));
    });
    }

    // const followorunfollow = document.getElementById("flw");
    // if (followorunfollow) {
    //     document.addEventListener("click", function(event) {

    //     })
    // }
}