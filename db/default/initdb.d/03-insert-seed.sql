USE `defaultdb`;

INSERT INTO `users` (
  id, display_name, icon_image_path, background_image_path, profile ,created_at, updated_at
)
VALUES
  ('DUMMY_USER_ID','tarou','icon_url','background_url','profile',NOW(),NOW());