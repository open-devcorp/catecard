// @ts-ignore
export let env = process.env.NODE_ENV || "development";

export let baseURL = env == "development" ? "http://localhost:3000" : "https://catecard.com";
