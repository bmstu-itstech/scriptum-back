package apiv2

import (
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
)

func valuesToDTO(vs []Value) []dto.Value {
	res := make([]dto.Value, len(vs))
	for i, v := range vs {
		res[i] = dto.Value{
			Type:  string(v.Type),
			Value: emptyOnNil(v.Value),
		}
	}
	return res
}

func valuesToAPI(vs []dto.Value) []Value {
	res := make([]Value, len(vs))
	for i, v := range vs {
		res[i] = Value{
			Type:  ValueType(v.Type),
			Value: nilOnEmpty(v.Value),
		}
	}
	return res
}

func fieldsToDTO(fs []Field) []dto.Field {
	res := make([]dto.Field, len(fs))
	for i, v := range fs {
		res[i] = dto.Field{
			Name: v.Name,
			Type: string(v.Type),
			Desc: nilOnNilOrEmpty(v.Desc),
			Unit: nilOnNilOrEmpty(v.Unit),
		}
	}
	return res
}

func fieldsToAPI(fs []dto.Field) []Field {
	res := make([]Field, len(fs))
	for i, v := range fs {
		res[i] = Field{
			Name: v.Name,
			Type: ValueType(v.Type),
		}
	}
	return res
}

func blueprintToAPI(b dto.BlueprintWithUser) Blueprint {
	return Blueprint{
		ArchiveID:  b.ArchiveID,
		CreatedAt:  b.CreatedAt,
		Desc:       nilOnNilOrEmpty(b.Desc),
		Id:         b.ID,
		In:         fieldsToAPI(b.In),
		Name:       b.Name,
		Out:        fieldsToAPI(b.Out),
		OwnerID:    b.OwnerID,
		OwnerName:  b.OwnerName,
		Visibility: Visibility(b.Visibility),
	}
}

func blueprintsToAPI(bs []dto.BlueprintWithUser) []Blueprint {
	res := make([]Blueprint, len(bs))
	for i, v := range bs {
		res[i] = blueprintToAPI(v)
	}
	return res
}

func jobToAPI(j dto.Job) Job {
	return Job{
		BlueprintID:   j.BlueprintID,
		BlueprintName: j.BlueprintName,
		CreatedAt:     j.CreatedAt,
		FinishedAt:    j.FinishedAt,
		Id:            j.ID,
		In:            fieldsToAPI(j.In),
		Input:         valuesToAPI(j.Input),
		Out:           fieldsToAPI(j.Out),
		Output:        valuesToAPI(j.Output),
		ResultCode:    j.ResultCode,
		ResultMsg:     j.ResultMsg,
		StartedAt:     j.StartedAt,
		State:         JobState(j.State),
	}
}

func jobsToAPI(js []dto.Job) []Job {
	res := make([]Job, len(js))
	for i, v := range js {
		res[i] = jobToAPI(v)
	}
	return res
}

func startJobRequestToDTO(r StartJobRequest, uid string, blueprintID string) request.StartJob {
	return request.StartJob{
		ActorID:     uid,
		BlueprintID: blueprintID,
		Values:      valuesToDTO(r.Values),
	}
}

func createBlueprintToDTO(r CreateBlueprintRequest, uid string) request.CreateBlueprint {
	return request.CreateBlueprint{
		ActorID:   uid,
		ArchiveID: r.ArchiveID,
		Name:      r.Name,
		Desc:      r.Desc,
		In:        fieldsToDTO(r.In),
		Out:       fieldsToDTO(r.Out),
	}
}

func userToAPI(u dto.User) User {
	return User{
		CreatedAt: u.CreatedAt,
		Email:     u.Email,
		Id:        u.ID,
		Name:      u.Name,
		Role:      Role(u.Role),
	}
}

func usersToAPI(us []dto.User) []User {
	res := make([]User, len(us))
	for i, u := range us {
		res[i] = userToAPI(u)
	}
	return res
}

func createUserToDTO(r CreateUserRequest, actorID string) request.CreateUser {
	return request.CreateUser{
		ActorID:  actorID,
		Email:    r.Email,
		Password: r.Password,
		Name:     r.Name,
		Role:     string(r.Role),
	}
}

func patchUserToDTO(r PatchUserRequest, actorID string, userID string) request.UpdateUser {
	return request.UpdateUser{
		ActorID:  actorID,
		UserID:   userID,
		Email:    r.Email,
		Password: r.Password,
		Name:     r.Name,
		Role:     r.Role,
	}
}
