import { makeAutoObservable } from "mobx";

class UserStore {
  isLoggedIn: boolean = false;
  username: string = '';

  constructor() {
    makeAutoObservable(this);
  }

  login(username: string) {
    this.isLoggedIn = true;
    this.username = username;
  }

  checkLogin(isLoggedIn: boolean, username:string) {
    this.isLoggedIn = isLoggedIn;
    this.username = username;
  }

  get isAdmin() {
    return this.username === 'admin';
  }

  logout() {
    this.isLoggedIn = false;
    this.username = '';
  }
}

export default new UserStore();