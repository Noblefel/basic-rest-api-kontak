ALTER TABLE public.users ADD COLUMN level INT DEFAULT(0);

INSERT INTO public.users (email, password, created_at, updated_at, level) VALUES
('admin@admin.com', '$2y$12$AI3YjckdxoalFHJQBGbQBu2aVNbaNpQO1wewFIaCrY5nMl4tnvYCq', '2020-01-02', '2020-01-02', 1);