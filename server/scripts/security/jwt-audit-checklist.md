# JWT Audit Checklist (Staging)

- Access token expiration <= 15 minutes.
- Refresh token rotation works and revocation on logout verified.
- Invalid signature tokens rejected.
- Expired tokens rejected.
- Role claims enforced on admin routes.
- No token leakage in logs.
- `Authorization` headers stripped from reverse proxy logs.
