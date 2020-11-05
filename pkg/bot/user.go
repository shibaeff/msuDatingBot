package bot

type User struct {
	name       string
	faculty    string
	gender     string
	wantGender string
	about      string
	id         int64
	photoLink  string
}

func (u *User) GetId() int64 {
	return u.id
}

func NewUser(
	name string,
	faculty string,
	gender string,
	wantGender string,
	about string,
	id int64,
	photoLink string,
) *User {
	return &User{
		name:       name,
		faculty:    faculty,
		gender:     gender,
		wantGender: wantGender,
		about:      about,
		id:         id,
		photoLink:  photoLink,
	}
}
