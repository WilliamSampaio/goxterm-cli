import { BACKEND_HOST } from "@/utils";
import axios from "axios";

export const getPing = () => {
  return axios.get(`http://${BACKEND_HOST}/api/ping`);
}

export const getInfo = () => {
  return axios.get(`http://${BACKEND_HOST}/api/info`);
}

export const getSshSessions = () => {
  return axios.get(`http://${BACKEND_HOST}/api/ssh/sessions`);
}
