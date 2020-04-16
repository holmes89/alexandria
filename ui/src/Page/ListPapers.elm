module Page.ListPapers exposing (Model, Msg, init, update, view)

import Html exposing (..)
import Html.Attributes exposing (class, href, src, style)
import Html.Events exposing (..)
import Http
import Paper exposing (Paper, papersDecoder)
import Session exposing (..)



-- MODEL


type Status
    = Failure
    | Loading
    | Success (List Paper)


type alias Model =
    { token : Token
    , status : Status
    }


init : Token -> ( Model, Cmd Msg )
init token =
    ( { token = token
      , status = Loading
      }
    , fetchPapers token
    )



-- UPDATE


type Msg
    = FetchPapers (Result Http.Error (List Paper))
    | Error


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        FetchPapers result ->
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
                    [ viewPapers model
                    ]
                ]
            ]
        ]


viewPapers : Model -> Html Msg
viewPapers model =
    case model.status of
        Failure ->
            div []
                [ text "Failed"
                ]

        Loading ->
            text "Loading..."

        Success papers ->
            div [ class "columns", class "is-mobile", class "is-multiline" ]
                (List.map
                    (\l ->
                        let
                            viewPath =
                                "/papers/" ++ l.id
                        in
                        div [ class "column", class "is-one-quarter" ]
                            [ div [ class "card" ]
                                [ header [ class "card-header" ]
                                    [ p [ class "card-header-title" ] [ text l.displayName ]
                                    ]
                                , div [ class "card-content", style "text-align" "center" ]
                                    [ img [ src ("http://read.jholmestech.com/assets/covers/" ++ l.id ++ ".jpg"), style "max-width" "300px", class "cover" ] []
                                    ]
                                , footer [ class "card-footer" ]
                                    [ a [ class "card-footer-item", href viewPath ] [ i [ class "fas", class "fa-paper-open" ] [], text "View" ]
                                    ]
                                ]
                            ]
                    )
                    papers
                )



-- HTTP


fetchPapers : Token -> Cmd Msg
fetchPapers token =
    Http.request
        { body = Http.emptyBody
        , expect = Http.expectJson FetchPapers papersDecoder
        , headers = [ Http.header "Authorization" token ]
        , method = "GET"
        , timeout = Nothing
        , tracker = Nothing
        , url = "https://docs.jholmestech.com/papers/"
        }
