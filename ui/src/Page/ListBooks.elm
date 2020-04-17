module Page.ListBooks exposing (Model, Msg, init, update, view)

import Book exposing (Book, booksDecoder)
import Dict exposing (Dict)
import Html exposing (..)
import Html.Attributes exposing (class, href, src, style)
import Html.Events exposing (..)
import Http
import Session exposing (..)
import Tag exposing (Tag, fetchTags)



-- MODEL


type Status
    = Failure
    | Loading
    | Success (List Book)


type alias Model =
    { token : Token
    , status : Status
    , tagDict : Dict String Tag
    }


init : Token -> ( Model, Cmd Msg )
init token =
    ( { token = token
      , status = Loading
      , tagDict = Dict.empty
      }
    , fetchTags token FetchTags
    )



-- UPDATE


type Msg
    = FetchBooks (Result Http.Error (List Book))
    | FetchTags (Result Http.Error (List Tag))
    | Error


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        FetchTags result ->
            case result of
                Ok tags ->
                    ( { model | tagDict = Dict.fromList (List.map (\e -> ( e.id, e )) tags) }, fetchBooks model.token )

                Err _ ->
                    ( { model | status = Failure }, Cmd.none )

        FetchBooks result ->
            case result of
                Ok url ->
                    ( { model | status = Success url }, Cmd.none )

                Err _ ->
                    ( { model | status = Failure }, Cmd.none )

        Error ->
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
                        div [ class "column", class "is-one-quarter-desktop is-full-mobile" ]
                            [ div [ class "card" ]
                                [ header [ class "card-header" ]
                                    [ p [ class "card-header-title" ] [ text l.displayName ]
                                    ]
                                , div [ class "card-content", style "text-align" "center" ]
                                    [ img [ src ("http://read.jholmestech.com/assets/covers/" ++ l.id ++ ".jpg"), style "max-width" "300px", class "cover" ] []
                                    , div [ class "tags" ]
                                        (List.map
                                            (\t ->
                                                let
                                                    tagEntry =
                                                        Dict.get t model.tagDict
                                                in
                                                case tagEntry of
                                                    Nothing ->
                                                        span [ class "tag" ] [ text "Unknown" ]

                                                    Just tag ->
                                                        span [ class "tag is-dark", style "background-color" tag.color ] [ text tag.displayName ]
                                            )
                                            l.tags
                                        )
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


fetchBooks : Token -> Cmd Msg
fetchBooks token =
    Http.request
        { body = Http.emptyBody
        , expect = Http.expectJson FetchBooks booksDecoder
        , headers = [ Http.header "Authorization" token ]
        , method = "GET"
        , timeout = Nothing
        , tracker = Nothing
        , url = "https://docs.jholmestech.com/books/"
        }
