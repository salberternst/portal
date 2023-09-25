package models

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
)

type OAuth2Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ThingsboardCredentials struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	Scope        string `json:"scope"`
}

type ThingsboardEntityId struct {
	Id string `json:"id"`
}

type ThingsboardCreateThing struct {
	Name  string `json:"name" binding:"required"`
	Label string `json:"label" binding:"required"`
	Model string `json:"model" binding:"required"`
}

type ThingsboardThing struct {
	Id    ThingsboardEntityId `json:"id"`
	Name  string              `json:"name"`
	Type  string              `json:"type"`
	Label string              `json:"label"`
}

type ThingsboardThingCredentials struct {
	CredentialsType string `json:"credentialsType"`
	CredentialsId   string `json:"credentialsId"`
}

type ThingsboardThingsList struct {
	Data []ThingsboardThing `json:"data"`
}

func GetThingsboardToken(accessToken string) (ThingsboardCredentials, error) {
	var thingsboardCredentials ThingsboardCredentials
	oauth2Credentials := OAuth2Credentials{
		Username: "oauth2-token",
		Password: accessToken,
	}

	err := requests.URL("http://thingsboard.192-168-178-60.nip.io/api/auth/login").
		BodyJSON(&oauth2Credentials).
		ToJSON(&thingsboardCredentials).
		Fetch(context.Background())

	if err != nil {
		return ThingsboardCredentials{}, nil
	}

	return thingsboardCredentials, nil
}

func GetThingsboardThings(accessToken string) (ThingsboardThingsList, error) {
	var thingsboardDeviceList ThingsboardThingsList

	err := requests.URL("http://thingsboard.192-168-178-60.nip.io/api/tenant/devices?page=0&pageSize=30").
		Header("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		ToJSON(&thingsboardDeviceList).
		Fetch(context.Background())

	fmt.Println(thingsboardDeviceList)

	if err != nil {
		return ThingsboardThingsList{}, err
	}

	return thingsboardDeviceList, nil
}

func CreateThingsboardThing(thing ThingsboardCreateThing, accessToken string) error {
	var thingsboardDevice ThingsboardThing

	err := requests.URL("http://thingsboard.192-168-178-60.nip.io/api/device").
		Header("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		BodyJSON(&thing).
		ToJSON(&thingsboardDevice).
		Fetch(context.Background())

	if err != nil {
		return err
	}

	return err
}

func GetThingsboardCredentials(id string, accessToken string) (ThingsboardThingCredentials, error) {
	var thingsboardCredentials ThingsboardThingCredentials

	err := requests.URL(fmt.Sprintf("http://thingsboard.192-168-178-60.nip.io/api/device/%s/credentials", id)).
		Header("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		ToJSON(&thingsboardCredentials).
		Fetch(context.Background())

	if err != nil {
		return ThingsboardThingCredentials{}, err
	}

	return thingsboardCredentials, nil
}

func DeleteThingsboardThing(id string, accessToken string) error {
	err := requests.URL(fmt.Sprintf("http://thingsboard.192-168-178-60.nip.io/api/device/%s", id)).
		Header("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		Delete().
		Fetch(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func GetThingsboardThing(id string, accessToken string) (ThingsboardThing, error) {
	var thingsboardThing ThingsboardThing

	err := requests.URL(fmt.Sprintf("http://thingsboard.192-168-178-60.nip.io/api/device/%s", id)).
		Header("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		ToJSON(&thingsboardThing).
		Fetch(context.Background())

	if err != nil {
		return ThingsboardThing{}, err
	}

	return thingsboardThing, nil
}
