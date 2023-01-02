let ws;
let receiver;
let username;
const host = "http://localhost:8800";

async function main() {
  const root = document.querySelector("#root");
  let res = await fetch(host + "/saas/api/auth", {
    method: "POST",
  });
  if (res.status === 200) {
    username = (await res.json()).data;
    console.log(username);
    renderChatApp();
    await renderUserList();
  } else {
    root.innerHTML = loginForm();
    document
      .querySelector(".login__submit")
      .addEventListener("click", async (e) => {
        e.preventDefault();
        let res = await fetch(host + "/saas/api/get-access-token", {
          method: "POST",
          body: JSON.stringify({
            username: document.querySelector("#username").value,
            password: document.querySelector("#password").value,
          }),
        });

        console.log(await res.json());
        location.reload();
      });
  }
}

main();

function renderChatApp(props) {
  ws = new WebSocket("ws://localhost:8800/saas/api/ws");
  ws.onmessage = (e) => {
    const message = JSON.parse(e.data);
    console.log(message);
    if (message.eventType === "Message") {
      if (message.sender == username) {
        renderMyMessage(message.eventPayload);
      } else {
        renderOtherMessage(message.eventPayload);
      }
    }
  };

  let chatApp = document.createElement("div");
  chatApp.className = "container";
  chatApp.innerHTML = `
  <div class="chat">
  <div class="row clearfix">
    <div class="col-lg-12">
      <div class="card chat-app">
        <div id="plist" class="people-list">
          <div class="input-group">
            <div class="input-group-prepend">
              <span class="input-group-text">
                <i class="fa fa-search"></i>
              </span>
            </div>
            <input
              type="text"
              class="form-control"
              placeholder="Search..."
            />
          </div>
          <ul class="list-unstyled chat-list mt-2 mb-0">
          
          </ul>
        </div>
        <div class="chat">
          <div class="chat-header clearfix">
            
          </div>
          <div class="chat-history">
            <ul class="message-list m-b-0">
              
            </ul>
          </div>
          <div class="chat-message clearfix">
            <div class="input-group mb-0">
              <button id="submit" class="btn input-group-prepend">
                <span class="input-group-text">
                  <i class="fa fa-send"></i>
                </span>
              </button>
              <input
                id="message"
                type="text"
                class="form-control"
                placeholder="Enter text here..."
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
  `;
  chatApp.querySelector("#submit").addEventListener("click", async (e) => {
    let chatRoomId;
    e.preventDefault();
    e.stopPropagation();
    const msg = document.getElementById("message").value;
    let res = await (
      await fetch(host + "/saas/api/chat-room/dm/members", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(receiver),
      })
    ).json();
    let chatRoom = res.data;
    console.log(chatRoom);
    if (chatRoom.members == null) {
      let result = await fetch(host + "/saas/api/chat-room/dm", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(receiver),
      });
      chatRoomId = (await result.json()).data.InsertedID;
    } else {
      chatRoomId = chatRoom.id;
    }
    console.log(chatRoomId);
    const message = {
      eventType: "Message",
      eventPayload: {
        chatroom: chatRoomId,
        content: msg,
      },
    };
    ws.send(JSON.stringify(message));
  });

  document.querySelector("#root").appendChild(chatApp);
}

function renderMyMessage(msg) {
  document.querySelector(".message-list").innerHTML += `
  
    <li class="clearfix">
      <div class="message-data text-left">
        <span class="message-data-time">10:16 AM, Today</span>
        <img
          src="https://bootdey.com/img/Content/avatar/avatar7.png"
          alt="avatar"
        />
      </div>
      <div class="message my-message float-right">
        ${msg.content}
      </div>
    </li>
  
  `;
}

function renderOtherMessage(msg) {
  document.querySelector(".message-list").innerHTML += `
    <li class="clearfix">
      <div class="message-data text-left">
        <span class="message-data-time">10:10 AM, Today</span>
        <img
          src="https://bootdey.com/img/Content/avatar/avatar7.png"
          alt="avatar"
        />
      </div>
      <div class="message other-message float-left">
        ${msg.content}
      </div>
    </li>
  `;
}

function loginForm() {
  return `
  <div class="container">
	<div class="screen">
		<div class="screen__content">
			<form class="login">
				<div class="login__field">
					<input id="username" type="text" class="login__input" placeholder="User name / Email">
				</div>
				<div class="login__field">
					<input id="password" type="password" class="login__input" placeholder="Password">
				</div>
				<button class="button login__submit">
					<span class="button__text">Log In Now</span>
				</button>				
			</form>
		</div>
		<div class="screen__background">
			<span class="screen__background__shape screen__background__shape4"></span>
			<span class="screen__background__shape screen__background__shape3"></span>		
			<span class="screen__background__shape screen__background__shape2"></span>
			<span class="screen__background__shape screen__background__shape1"></span>
		</div>		
	</div>
</div>
  `;
}

async function renderUserList() {
  let res = await fetch(host + "/saas/api/user/online");
  body = await res.json();
  console.log(body);
  let result = "";
  document.querySelector(".chat-list").innerHTML = "";
  for (let i = 0; i < body.data.length; i++) {
    if (body.data[i] == username) {
      continue;
    }
    let user = document.createElement("li");
    user.classList.add("user", "clearfix");
    user.innerHTML = `
    <div class="about">
    <div class="name">${body.data[i]}</div>
    <div class="status">
      <i class="fa fa-circle online"></i> online
    </div>
  </div>
  `;
    user.addEventListener("click", async () => {
      renderChatHeader(body.data[i]);
    });
    document.querySelector(".chat-list").appendChild(user);
  }
}

function renderChatHeader(_receiver) {
  receiver = _receiver;
  document.querySelector(".chat-header").innerHTML = `
  <div class="chat-message clearfix">
                    <div class="col-lg-6">
                      <a
                        href="javascript:void(0);"
                        data-toggle="modal"
                        data-target="#view_info"
                      >
                        <img
                          src="https://bootdey.com/img/Content/avatar/avatar2.png"
                          alt="avatar"
                        />
                      </a>
                      <div class="chat-about">
                        <h6 class="m-b-0">${_receiver}</h6>
                      </div>
                    </div>
                  </div>
  `;
}

async function renderChatHistory(chatroom) {
  let res = await fetch(host + "/saas/api/message/chat-room/" + chatroom);
  body = await res.json();
  for (let i = 0; i < body.data; i++) {
    if (body.data[i].sender == username) {
      renderMyMessage(body.data[i]);
    } else {
      renderOtherMessage(body.data[i]);
    }
  }
}

// async function renderMessages() {
//   let res = await fetch("http://localhost:8800/api/chat-room/members", {
//     method: "POST",
//     headers: {

//       "Content-Type": "application/json",
//     }
//     body: JSON.stringify({[""]})
//   })

//   res =await fetch("http://localhost:8800/api/message/pagination", {
//     method: "POST",
//     headers: {
//       "Content-Type": "application/json",
//     },
//     body: JSON.stringify({
//       chatRoomId: "",
//       limit: 10,
//       skip: 0,
//     }),
//   });
// }
