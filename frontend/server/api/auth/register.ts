/**
 * Server-side handler for the `/api/auth/register` route that sends a POST request to the server to register a new user.
 */
export default defineEventHandler(async (event) => {
  const body = await readBody(event);

  const response = await $fetch("http://server:8080/api/v1/user/register", {
    method: "POST",
    body: {
      email: body.email,
      username: body.username,
      password: body.password,
    },
  });

  return response;
});
