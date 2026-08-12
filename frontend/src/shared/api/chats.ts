import { apiRequest } from "./client";

export default {
    getList: async () => apiRequest("/me/chats"),
};