import React, { useEffect } from "react";
import {
  BrowserRouter as Router,
  Route,
  NavLink,
  Redirect,
} from "react-router-dom";
import { observer } from "mobx-react-lite";
import { Breadcrumb, Button } from "@arco-design/web-react";
import UserStore from "./store";
import Login from "./components/Login";
import Register from "./components/Register";
import CandidateList from "./components/CandidateList";
import VoteResult from "./components/VoteResult";
import "@arco-design/web-react/dist/css/arco.css";
import { logout } from "./services/api";
import Vote from "./components/Vote";

const routesName: any = {
  register: "注册",
  login: "登录",
  "vote-candidate": "投票",
  candidates: "候选人列表",
  "vote-result": "投票结果",
};

const App: React.FC = observer(() => {
  useEffect(() => {
    //模拟已经登录
   //UserStore.checkLogin(true, "testuser");

    // 模拟管理员登录
    //UserStore.checkLogin(true, "admin");
  });
  console.log(UserStore.username, UserStore.isLoggedIn);
  const handleLogout = async () => {
    try {
      await logout();
    } catch (error) {
      console.log(error);
    }

    window.location.replace("/");
  };
  return (
    <Router>
      <div style={{ maxWidth: "1000px", margin: "0 auto" }}>
        <Breadcrumb style={{ margin: "16px 0" }}>
          <Breadcrumb.Item>
            <NavLink to="/">首页</NavLink>
          </Breadcrumb.Item>
          <Breadcrumb.Item>
            <Route
              path="/"
              render={({ location }) => {
                const routes = location.pathname.split("/").filter(Boolean);
                const currentRoute = routes[routes.length - 1];
                return <span>{routesName[currentRoute]}</span>;
              }}
            />
          </Breadcrumb.Item>
        </Breadcrumb>

        <Route exact path="/">
          <div style={{ textAlign: "center" }}>
            <h1 style={{ fontSize: "20px" }}>欢迎使用匿名投票系统</h1>
            <div style={{ margin: "0 auto" }}>
              {!UserStore.isLoggedIn ? (
                <>
                  <Button
                    type="primary"
                    style={{
                      display: "block",
                      width: "100px",
                      margin: "0 auto",
                      marginBottom: "8px",
                    }}
                  >
                    <NavLink to="/register">注册</NavLink>
                  </Button>
                  <Button
                    type="primary"
                    style={{
                      display: "block",
                      width: "100px",
                      margin: "0 auto",
                      marginBottom: "8px",
                    }}
                  >
                    <NavLink to="/login">登录</NavLink>
                  </Button>
                </>
              ) : (
                <>
                  <Button
                    type="primary"
                    style={{
                      display: "block",
                      width: "100px",
                      margin: "0 auto",
                      marginBottom: "8px",
                    }}
                  >
                    <NavLink to="/vote-candidate">投票</NavLink>
                  </Button>
                  <Button
                    type="primary"
                    onClick={handleLogout}
                    style={{
                      display: "block",
                      width: "100px",
                      margin: "0 auto",
                      marginBottom: "8px",
                    }}
                  >
                    登出
                  </Button>
                </>
              )}
              {UserStore.isLoggedIn && UserStore.isAdmin ? (
                <>
                  <Button
                    type="primary"
                    style={{
                      display: "block",
                      width: "100px",
                      margin: "0 auto",
                      marginBottom: "8px",
                    }}
                  >
                    <NavLink to="/candidates">查看候选人</NavLink>
                  </Button>
                </>
              ) : null}

              <Button
                type="primary"
                style={{
                  display: "block",
                  width: "140px",
                  margin: "0 auto",
                  marginBottom: "8px",
                }}
              >
                <NavLink to="/vote-result">查看投票结果</NavLink>
              </Button>
            </div>
          </div>
        </Route>

        <Route path="/register">
          {!UserStore.isLoggedIn ? <Register /> : <Redirect to="/" />}
        </Route>

        <Route path="/login">
          {!UserStore.isLoggedIn ? <Login /> : <Redirect to="/" />}
        </Route>

        <Route path="/candidates">
          {UserStore.isLoggedIn && UserStore.isAdmin ? (
            <CandidateList />
          ) : (
            <Redirect to="/" />
          )}
        </Route>

        <Route path="/vote-result">
          {UserStore.isLoggedIn ? <VoteResult /> : <Redirect to="/" />}
        </Route>
        <Route path="/vote-candidate">
          {UserStore.isLoggedIn ? <Vote /> : <Redirect to="/" />}
        </Route>
      </div>
    </Router>
  );
});

export default App;
