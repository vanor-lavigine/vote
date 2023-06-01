import axios from "axios";
import UserStore from "../store";
import { Message } from "@arco-design/web-react";

axios.defaults.baseURL = "http://localhost:8080";
axios.defaults.headers.post["Content-Type"] = "application/json";

type ApiResponse = {
  code: number;
  data: any;
  errorMessage?: string;
  errorcode?: string;
};

axios.interceptors.response.use(
  function (response) {
    if (response?.data?.code !== 200) {
      if (response?.data?.ErrorMessage) {
        Message.error(response?.data?.ErrorMessage);
      }
    }
    return response.data;
  },
  function (error) {
    // return Promise.reject(error);
    console.log(error);
    Message.error("网络错误");
  }
);

export async function login(username: string, password: string): Promise<void> {
  const data:any = await axios.post("/login", { username, password });
  if (data?.code === 200) {
    UserStore.login(username);
    //window.location.pathname = '/';
  }
}

export async function register(
  username: string,
  password: string
): Promise<void> {
  const data:any = await axios.post("/register", { username, password });
  if (data?.code === 200) {
    await login(username, password);
  }
}

export async function listCandidates(): Promise<ApiResponse> {
  const data:any = await axios.post("/listCandidates");
  //console.log(data)
  return data?.data;

}

export async function deleteCandidate(id: number): Promise<ApiResponse> {
  const data:any = await axios.post("/deleteCandidate", { id });
  console.log(data)
  return data?.data;
}

export async function createCandidate(username: string): Promise<ApiResponse> {
  const data:any = await axios.post("/createCandidate", { username });
  return data?.data;
}

export async function getVoteList(): Promise<ApiResponse> {
  const data:any = await axios.post("/getVoteList");
  //console.log(data)
  return data?.data;
}

export async function logout(): Promise<void> {
  try {
    await axios.post("/logout");
  } catch (error) {
    console.log("logout error", error);
  }
  UserStore.logout();
}

export async function checkLogin(): Promise<void> {
  const data:any = await axios.post("/checkLogin");
  if (data?.code === 200) {
    const { HasLogin, UserName } = data?.data || {};
    UserStore.checkLogin(HasLogin, UserName);
  }
}
