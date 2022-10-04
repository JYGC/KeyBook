CREATE OR REPLACE FUNCTION
	"KeyBook".sp_GetDeviceActivityHistory
	(
		DeviceId uuid
	)
RETURNS TABLE (
	RecordDateTime timestamp with time zone,
	Description text
) AS
	$func$
BEGIN
	RETURN QUERY
	SELECT
		historyList."RecordDateTime", historyList."Description"
	FROM
	(
		SELECT
			dh."DeviceId", dh."RecordDateTime", dh."Description"
		FROM
			"KeyBook"."DeviceHistory" dh
		WHERE
			dh."DeviceId" = DeviceId

		UNION

		SELECT
			pdh."DeviceId", pdh."RecordDateTime", pdh."Description"
		FROM
			"KeyBook"."PersonDeviceHistory" pdh
		WHERE
			pdh."DeviceId" = DeviceId
	) historyList
	ORDER BY
		historyList."DeviceId" DESC,
		historyList."RecordDateTime" DESC;
END
$func$ LANGUAGE plpgsql;