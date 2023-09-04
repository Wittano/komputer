import com.github.jengelman.gradle.plugins.shadow.tasks.ShadowJar
import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    kotlin("jvm") version "1.9.10"
    kotlin("kapt") version "1.9.10"
    id("com.github.johnrengelman.shadow") version "8.1.1"

    application
}

group = "com.wittano.komputer.cli"
version = rootProject.version

repositories {
    mavenCentral()
}

val picocliVersion = "4.7.4"

dependencies {
    // Internal dependencies
    implementation(project(":core"))

    // Kotlin
    implementation(kotlin("stdlib"))

    // Discord4j
    // TODO Set global version on used dependencies
    implementation("com.discord4j:discord4j-core:3.2.5")
    implementation("io.projectreactor.kotlin:reactor-kotlin-extensions:1.2.2")

    // Picocli
    implementation("info.picocli:picocli:$picocliVersion")
    kapt("info.picocli:picocli-codegen:$picocliVersion")

    // Logger
    implementation("ch.qos.logback:logback-classic:1.4.11")
    implementation("org.codehaus.janino:janino:3.1.10")

    testImplementation(platform("org.junit:junit-bom:5.9.1"))
    testImplementation("org.junit.jupiter:junit-jupiter")
}

// TODO Create optional native image with GrallVM
// TODO Generate script to run CLI

tasks.test {
    useJUnitPlatform()
}

application {
    mainClass.set("com.wittano.komputer.cli.MainKt")
}

tasks.withType<ShadowJar> {
    manifest {
        attributes["Main-Class"] = application.mainClass
    }

    archiveBaseName.set("komputer-cli")
    archiveClassifier.set("")
}

tasks.withType<KotlinCompile> {
    kotlinOptions.jvmTarget = "17"
}