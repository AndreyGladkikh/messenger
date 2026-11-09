import { apiRequest } from "./client";

export default {
    getList: async () => apiRequest("/bff/web/home-page"),
};
