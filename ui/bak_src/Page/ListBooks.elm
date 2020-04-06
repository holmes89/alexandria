module Page.ListBooks exposing (Model, Msg, init, update, view)

import Book exposing (Book, booksDecoder)
import Html exposing (..)
import Html.Attributes exposing (class, href, src, style)
import Html.Events exposing (..)
import Http
import Session exposing (..)



-- MODEL


type Status
    = Failure
    | Loading
    | Success (List Book)


type alias Model =
    { session : Session
    , status : Status
    }


init : Session -> ( Model, Cmd Msg )
init session =
    ( { session = session
      , status = Loading
      }
    , fetchBooks
    )



-- UPDATE


type Msg
    = FetchBooks (Result Http.Error (List Book))


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        FetchBooks result ->
            case result of
                Ok url ->
                    ( { model | status = Success url }, Cmd.none )

                Err _ ->
                    ( { model | status = Failure }, Cmd.none )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- VIEW


view : Model -> Html Msg
view model =
    div []
        [ div []
            [ section [ class "section" ]
                [ div [ class "container" ]
                    [ viewBooks model
                    ]
                ]
            ]
        ]


viewBooks : Model -> Html Msg
viewBooks model =
    case model.status of
        Failure ->
            div []
                [ text "Failed"
                ]

        Loading ->
            text "Loading..."

        Success books ->
            div [ class "columns", class "is-mobile", class "is-multiline" ]
                (List.map
                    (\l ->
                        let
                            viewPath =
                                "/books/" ++ l.id
                        in
                        div [ class "column", class "is-one-quarter" ]
                            [ div [ class "card" ]
                                [ header [ class "card-header" ]
                                    [ p [ class "card-header-title" ] [ text l.displayName ]
                                    ]
                                , div [ class "card-content", style "text-align" "center" ]
                                    [ img [ src ("http://read.jholmestech.com/assets/covers/" ++ l.id ++ ".jpg"), style "max-width" "300px" ] []
                                    ]
                                , footer [ class "card-footer" ]
                                    [ a [ class "card-footer-item", href viewPath ] [ i [ class "fas", class "fa-book-open" ] [], text "View" ]
                                    ]
                                ]
                            ]
                    )
                    books
                )



-- HTTP


fetchBooks : Cmd Msg
fetchBooks =
    Http.get
        { url = "https://docs.jholmestech.com/books/"
        , expect = Http.expectJson FetchBooks booksDecoder
        }
