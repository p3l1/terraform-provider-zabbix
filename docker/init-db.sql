-- ABOUTME: Database initialization script for test environment.
-- ABOUTME: Creates a static API token for acceptance tests.

-- Insert API token for CI/CD testing
-- Token value: 071fb9d2e8f72cf9c40128f0f5aab3def1bab0893413314b083fdcb4551eb01a
INSERT INTO public.token (tokenid, name, description, userid, token, lastaccess, status, expires_at, created_at, creator_userid)
VALUES (1, 'cicd', 'CI/CD', 1, '4b99fc440dabdd5b3dca02b16c1a5a705f01fd25d1a8f7ef0743cd9029499ebe09d179bba2bd3408919de75d0f60330b323d7838b65f979e5bb59382577d7a83', 0, 0, 0, 1766783844, 1)
ON CONFLICT (tokenid) DO NOTHING;
