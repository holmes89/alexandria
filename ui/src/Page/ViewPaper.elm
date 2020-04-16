module Page.ViewPaper exposing (Model, Msg, init, update, view)

import Browser.Navigation as Nav
import Html exposing (..)
import Html.Attributes exposing (class, href, src, style)
import Http
import Paper exposing (Paper, PaperID, paperDecoder)
import Session exposing (..)


type alias Model =
    { navKey : Nav.Key
    , status : Status
    , token : Token
    }


type Status
    = Failure
    | Loading
    | Success Paper


init : PaperID -> Nav.Key -> Token -> ( Model, Cmd Msg )
init paperID navKey token =
    ( initialModel navKey token, getPaper paperID token )


initialModel : Nav.Key -> Token -> Model
initialModel navKey token =
    { navKey = navKey
    , status = Loading
    , token = token
    }


type Msg
    = FetchPaper (Result Http.Error Paper)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        FetchPaper result ->
            case result of
                Ok url ->
                    ( { model | status = Success url }, Cmd.none )

                Err _ ->
                    ( { model | status = Failure }, Cmd.none )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none


view : Model -> Html Msg
view model =
    case model.status of
        Failure ->
            div []
                [ text "Failed"
                ]

        Loading ->
            text "Loading..."

        Success paper ->
            div []
                [ section [ class "hero is-light" ]
                    [ div [ class "hero-body" ]
                        [ div [ class "container" ]
                            [ h1 [ class "title" ] [ text paper.displayName ]
                            , h2 [ class "subtitle" ] [ text paper.description ]
                            ]
                        ]
                    ]
                , section [ class "section" ]
                    [ div [ class "container" ]
                        [ div [ class "columns is-centered is-mobile" ]
                            [ div [ class "column", class "is-4" ]
                                [ aside [ class "menu" ]
                                    [ p [ class "menu-label" ] [ text "Options" ]
                                    , ul [ class "menu-list" ]
                                        [ li []
                                            [ a [] [ text "Edit" ]
                                            , a [ href paper.path ] [ text "Download" ]
                                            ]
                                        ]
                                    ]
                                ]
                            , div [ class "column is-4" ]
                                [ img [ src ("http://read.jholmestech.com/assets/covers/" ++ paper.id ++ ".jpg"), style "max-width" "300px" ] [] ]
                            ]
                        ]
                    ]
                ]



-- HTTP


getPaper : PaperID -> Token -> Cmd Msg
getPaper paperID token =
    Http.request
        { body = Http.emptyBody
        , expect = Http.expectJson FetchPaper paperDecoder
        , headers = [ Http.header "Authorization" token ]
        , method = "GET"
        , timeout = Nothing
        , tracker = Nothing
        , url = "https://docs.jholmestech.com/documents/" ++ paperID
        }
