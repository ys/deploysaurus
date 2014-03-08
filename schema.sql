CREATE EXTENSION "uuid-ossp";

CREATE TABLE users(
  id uuid primary key default uuid_generate_v4(),
  email text,
  github_login text,
  github_id text,
  github_token text,
  heroku_id text,
  heroku_token text,
  heroku_refresh_token text,
  heroku_expiration timestamptz
);

CREATE INDEX index_on_users_for_github_login ON users(github_login);
CREATE INDEX index_on_users_for_heroku_id ON users(heroku_id);
