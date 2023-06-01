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
  const response = await axios.post("/login", { username, password });
  if (response?.data?.code === 200) {
    UserStore.login(username);
  }
}

export async function register(
  username: string,
  password: string
): Promise<void> {
  const response = await axios.post("/register", { username, password });
  if (response?.data?.code === 200) {
    await login(username, password);
  }
}

export async function listCandidates(): Promise<ApiResponse> {
  const response = await axios.post("/listCandidates");
  return response?.data?.Data;

}

export async function deleteCandidate(id: number): Promise<ApiResponse> {
  const response = await axios.post("/deleteCandidate", { id });
  return response?.data?.Data;
}

export async function createCandidate(username: string): Promise<ApiResponse> {
  const response = await axios.post("/createCandidate", { username });
  return response?.data?.Data;
}

export async function getVoteList(): Promise<ApiResponse> {
  const response = await axios.post("/getVoteList");
  return response?.data?.Data;
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
  const response = await axios.post("/checkLogin");
  if (response?.data?.code === 200) {
    const { HasLogin, UserName } = response?.data?.Data || {};
    UserStore.checkLogin(HasLogin, UserName);
  }
}
