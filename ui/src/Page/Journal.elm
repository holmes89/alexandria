module Page.Journal exposing (Model, Msg, init, update, view)

import Html exposing (..)
import Html.Attributes exposing (class, href, id, rows, src, style)
import Http
import Journal exposing (Entry, entriesDecoder)
import Session exposing (..)



-- MODEL


type Status
    = Failure
    | Loading
    | Success (List Entry)


type alias Model =
    { token : Token
    , status : Status
    }


init : Token -> ( Model, Cmd Msg )
init token =
    ( { token = token
      , status = Loading
      }
    , fetchEntries token
    )



-- UPDATE


type Msg
    = FetchEntries (Result Http.Error (List Entry))
    | Error


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        FetchEntries result ->
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


view : Model -> Html msg
view model =
    case model.status of
        Failure ->
            div []
                [ text "Failed"
                ]

        Loading ->
            text "Loading..."

        Success entries ->
            section [ class "hero" ]
                [ div [ class "hero-body journal-entry-area" ]
                    [ div [ class "container" ]
                        (List.map
                            viewEntry
                            entries
                        )
                    ]
                , div [ class "hero-foot journal-draft-area" ]
                    [ div [ class "container" ]
                        [ div [ class "columns" ]
                            [ div [ class "column is-11" ]
                                [ textarea [ class "textarea", rows 5 ] [] ]
                            , div [ class "column journal-entry-button" ]
                                [ button [ class "button" ] [ text "Submit" ] ]
                            ]
                        ]
                    ]
                ]


viewEntry : Entry -> Html msg
viewEntry entry =
    div [ class "columns" ]
        [ div [ class "column is-1" ] [ text entry.created ]
        , div [ class "column" ]
            [ article [ class "message is-dark", id entry.id ]
                [ div [ class "message-body" ] [ text entry.content ]
                ]
            ]
        ]



-- HTTP


fetchEntries : Token -> Cmd Msg
fetchEntries token =
    Http.request
        { body = Http.emptyBody
        , expect = Http.expectJson FetchEntries entriesDecoder
        , headers = [ Http.header "Authorization" token ]
        , method = "GET"
        , timeout = Nothing
        , tracker = Nothing
        , url = "https://docs.jholmestech.com/journal/entry/"
        }
