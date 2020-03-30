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


view : Model -> Html Msg
view model =
    div []
        [ div []
            [ section [ class "section" ]
                [ div [ class "container" ] [ text "here" ]
                ]
            ]
        ]



-- HTTP


getBook : BookID -> Cmd Msg
getBook bookID =
    Http.get
        { url = "https://docs.jholmestech.com/books/" ++ bookID
        , expect = Http.expectJson FetchBook bookDecoder
        }
