CREATE OR REPLACE FUNCTION
	"KeyBook".sp_GetDeviceActivityHistory
	(
		DeviceId uuid,
		OrganizationId uuid
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
			"KeyBook"."Device" d
		INNER JOIN
			"KeyBook"."DeviceHistory" dh ON d."Id" = dh."DeviceId"
		WHERE
			d."Id" = DeviceId
		AND
			d."OrganizationId" = OrganizationId

		UNION

		SELECT
			pdh."DeviceId", pdh."RecordDateTime", pdh."Description"
		FROM
			"KeyBook"."Device" d
		INNER JOIN
			"KeyBook"."PersonDeviceHistory" pdh ON d."Id" = pdh."DeviceId"
		WHERE
			d."Id" = DeviceId
		AND
			d."OrganizationId" = OrganizationId
	) historyList
	ORDER BY
		historyList."DeviceId" DESC,
		historyList."RecordDateTime" DESC;
END
$func$ LANGUAGE plpgsql;