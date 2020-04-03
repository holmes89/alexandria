module Page.Login exposing (..)

import Base64
import Browser.Navigation as Nav
import Html exposing (..)
import Html.Attributes exposing (action, attribute, class, for, href, placeholder, required, src, style, type_)
import Html.Events exposing (onClick, onInput)
import Http
import Json.Decode as Decode exposing (Decoder, field, string)
import Json.Decode.Pipeline as DecodePipeline
import Session exposing (..)
import Url exposing (Url)



-- MODEL


type alias Model =
    { session : Session
    , navKey : Nav.Key
    , username : String
    , password : String
    }


type alias Token =
    { token : String
    }


tokenDecoder : Decoder Token
tokenDecoder =
    Decode.succeed Token
        |> DecodePipeline.required "token" string


type Msg
    = UpdateUsername String
    | UpdatePassword String
    | SubmitLogin
    | Login (Result Http.Error Token)


init : Nav.Key -> ( Model, Cmd Msg )
init navKey =
    ( { session = Unauthenticated
      , navKey = navKey
      , username = ""
      , password = ""
      }
    , Cmd.none
    )


updateUsername : String -> Model -> Model
updateUsername username model =
    { model | username = username }


updatePassword : String -> Model -> Model
updatePassword password model =
    { model | password = password }


submitLogin : Model -> Cmd Msg
submitLogin model =
    let
        auth =
            "Basic " ++ Base64.encode (String.trim model.username ++ ":" ++ String.trim model.password)
    in
    Http.request
        { body = Http.emptyBody
        , expect = Http.expectJson Login tokenDecoder
        , headers = [ Http.header "Authorization" auth ]
        , method = "GET"
        , timeout = Nothing
        , tracker = Nothing
        , url = "https://docs.jholmestech.com/auth/"
        }


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UpdateUsername username ->
            ( updateUsername username model, Cmd.none )

        UpdatePassword password ->
            ( updatePassword password model, Cmd.none )

        SubmitLogin ->
            ( model, submitLogin model )

        Login result ->
            case result of
                Ok url ->
                    ( { model | session = Authenticated url.token }, Nav.pushUrl model.navKey "/books" )

                Err _ ->
                    ( { model | session = Unauthenticated }, Cmd.none )


view : Model -> Html Msg
view model =
    section [ class "hero is-light is-fullheight" ]
        [ div [ class "hero-body" ]
            [ div [ class "container" ]
                [ div [ class "columns is-centered" ]
                    [ div [ class "column is-5-tablet is-4-desktop is-3-widescreen" ]
                        [ h1 [ class "title" ] [ text "Login" ]
                        , div [ class "box" ]
                            [ div [ class "field" ]
                                [ label [ class "label", for "" ]
                                    [ text "Email" ]
                                , div [ class "control has-icons-left" ]
                                    [ input [ class "input", placeholder "e.g. bobsmith@gmail.com", attribute "required" "", type_ "email", onInput UpdateUsername ]
                                        []
                                    , span [ class "icon is-small is-left" ]
                                        [ i [ class "fa fa-envelope" ]
                                            []
                                        ]
                                    ]
                                ]
                            , div [ class "field" ]
                                [ label [ class "label", for "" ]
                                    [ text "Password" ]
                                , div [ class "control has-icons-left" ]
                                    [ input [ class "input", placeholder "*******", attribute "required" "", type_ "password", onInput UpdatePassword ]
                                        []
                                    , span [ class "icon is-small is-left" ]
                                        [ i [ class "fa fa-lock" ]
                                            []
                                        ]
                                    ]
                                ]
                            , div [ class "field" ]
                                [ button [ class "button is-link", onClick SubmitLogin ]
                                    [ text "Submit" ]
                                ]
                            ]
                        ]
                    ]
                ]
            ]
        ]
