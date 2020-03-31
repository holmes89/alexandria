module Page.ViewBook exposing (Model, Msg, init, update, view)

import Book exposing (Book, BookID, bookDecoder)
import Browser.Navigation as Nav
import Html exposing (..)
import Html.Attributes exposing (class, href, src, style)
import Http


type alias Model =
    { navKey : Nav.Key
    , status : Status
    }


type Status
    = Failure
    | Loading
    | Success Book


init : BookID -> Nav.Key -> ( Model, Cmd Msg )
init bookID navKey =
    ( initialModel navKey, getBook bookID )


initialModel : Nav.Key -> Model
initialModel navKey =
    { navKey = navKey
    , status = Loading
    }


type Msg
    = FetchBook (Result Http.Error Book)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        FetchBook result ->
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

        Success book ->
            div []
                [ section [ class "hero is-link" ]
                    [ div [ class "hero-body" ]
                        [ div [ class "container" ]
                            [ h1 [ class "title" ] [ text book.displayName ]
                            , h2 [ class "subtitle" ] [ text book.description ]
                            ]
                        ]
                    ]
                , section [ class "section" ]
                    [ div [ class "container" ]
                        [ div [ class "columns", class "is-mobile" ]
                            [ div [ class "column", class "is-3", class "is-offset-3" ]
                                [ aside [ class "menu" ]
                                    [ p [ class "menu-label" ] [ text "Options" ]
                                    , ul [ class "menu-list" ]
                                        [ li []
                                            [ a [] [ text "Edit" ]
                                            , a [ href book.path ] [ text "Download" ]
                                            ]
                                        ]
                                    ]
                                ]
                            , div [ class "column is-3" ]
                                [ img [ src ("http://read.jholmestech.com/assets/covers/" ++ book.id ++ ".jpg"), style "max-width" "300px" ] [] ]
                            ]
                        ]
                    ]
                ]



-- HTTP


getBook : BookID -> Cmd Msg
getBook bookID =
    Http.get
        { url = "https://docs.jholmestech.com/documents/" ++ bookID
        , expect = Http.expectJson FetchBook bookDecoder
        }
