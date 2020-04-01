module Page.Login exposing (..)

import Html exposing (..)
import Html.Attributes exposing (action, attribute, class, for, href, placeholder, required, src, style, type_)
import Html.Events exposing (onClick, onInput)
import Http
import Json.Decode as Decode exposing (Decoder, field, string)
import Json.Decode.Pipeline as DecodePipeline



-- MODEL


type alias Model =
    { username : String
    , password : String
    , token : String
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


init : ( Model, Cmd Msg )
init =
    ( { username = ""
      , password = ""
      , token = ""
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
    Http.get
        { url = "https://" ++ model.username ++ ":" ++ model.password ++ "@docs.jholmestech.com/auth/"
        , expect = Http.expectJson Login tokenDecoder
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
                    ( { model | token = url.token }, Cmd.none )

                Err _ ->
                    ( { model | token = "invalid login" }, Cmd.none )


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
        , div [] [ text model.token ]
        ]
