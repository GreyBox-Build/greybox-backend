INSERT INTO
    kyc_data (user_id, front_photo, back_photo)
SELECT
    user_id,
    front_photo,
    back_photo
FROM
    kycs;

ALTER TABLE kycs
DROP COLUMN front_photo,
DROP COLUMN back_photo;
