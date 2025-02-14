package pg

const insertUserQuery = `
	INSERT INTO posts.user (
		name
		, role
	) VALUES (
		$1, $2
	) RETURNING 
		id
		, name
		, role
		, created_at
`

const insertSessionQuery = `
	INSERT INTO posts.session (
		user_id
		, expires_at
	) VALUES (
		$1, $2 
	) RETURNING
	 	id
		, user_id
		, created_at
		, expires_at
`

const getSessionById = `
	SELECT 
		*
	FROM 
		posts.session
	WHERE
		id=$1
`

const getUserIdQuery = `
	SELECT 
		id
	FROM
		posts.user
	WHERE
		name=$1
`

const deleteSessionById = `
	DELETE FROM
		posts.session
	WHERE
		id=$1
`
